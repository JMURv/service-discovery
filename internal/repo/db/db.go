package db

import (
	"context"
	"errors"
	"github.com/JMURv/service-discovery/internal/repo"
	md "github.com/JMURv/service-discovery/pkg/model"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type Repository struct {
	mu      sync.Mutex
	conn    *gorm.DB
	rrIndex map[string]int
}

func New() *Repository {
	conn, err := gorm.Open(sqlite.Open("discovery.db"), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("failed to connect to the database", zap.Error(err))
	}
	if err = conn.AutoMigrate(&md.Service{}); err != nil {
		zap.L().Fatal("failed to migrate the database", zap.Error(err))
	}

	return &Repository{
		conn:    conn,
		rrIndex: make(map[string]int),
	}
}

func (r *Repository) Close() error {
	db, err := r.conn.DB()
	if err != nil {
		zap.L().Error("failed to get the database", zap.Error(err))
		return err
	}

	if err = db.Close(); err != nil {
		zap.L().Error("failed to close the database", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) Register(ctx context.Context, name, addr string) error {
	var svc md.Service

	if err := r.conn.WithContext(ctx).
		Where("name = ? AND address = ?", name, addr).
		First(&svc).Error; err == nil {
		return repo.ErrAlreadyExists
	}

	service := md.Service{Name: name, Address: addr}
	if err := r.conn.WithContext(ctx).Create(&service).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) Deregister(ctx context.Context, name, addr string) error {
	var service md.Service

	if err := r.conn.WithContext(ctx).
		Where("name = ? AND address = ?", name, addr).
		First(&service).Error; err != nil {
		return repo.ErrNotFound
	}

	if err := r.conn.WithContext(ctx).Delete(&service).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindServiceByName(ctx context.Context, name string) (string, error) {
	var svcs []md.Service

	if err := r.conn.WithContext(ctx).
		Where("name = ? AND is_active = true", name).
		Find(&svcs).Error; err != nil || len(svcs) == 0 {
		return "", repo.ErrNotFound
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	currentIndex := r.rrIndex[name]
	selectedAddr := svcs[currentIndex].Address

	r.rrIndex[name] = (currentIndex + 1) % len(svcs)
	return selectedAddr, nil
}

func (r *Repository) ListServices(ctx context.Context) ([]string, error) {
	var svcs []md.Service
	if err := r.conn.WithContext(ctx).
		Find(&svcs).Error; err != nil {
		return nil, err
	}

	names := make([]string, len(svcs))
	for i, svc := range svcs {
		names[i] = svc.Name
	}

	return names, nil
}

func (r *Repository) ListAddrs(ctx context.Context, name string) ([]string, error) {
	var svcs []md.Service
	if err := r.conn.WithContext(ctx).
		Where("name = ?", name).
		Find(&svcs).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	addrs := make([]string, len(svcs))
	for i, svc := range svcs {
		addrs[i] = svc.Address
	}

	return addrs, nil
}

func (r *Repository) DeactivateSvc(ctx context.Context, name, addr string) error {
	var svc md.Service

	if err := r.conn.WithContext(ctx).
		Where("name = ? AND address = ?", name, addr).
		First(&svc).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	svc.IsActive = false
	if err := r.conn.WithContext(ctx).Save(&svc).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) ActivateSvc(ctx context.Context, name, addr string) error {
	var svc md.Service

	if err := r.conn.WithContext(ctx).
		Where("name = ? AND address = ?", name, addr).
		First(&svc).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	svc.IsActive = true
	if err := r.conn.WithContext(ctx).Save(&svc).Error; err != nil {
		return err
	}
	return nil
}
