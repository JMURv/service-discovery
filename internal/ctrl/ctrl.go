package ctrl

import (
	"context"
	"errors"
	"github.com/JMURv/service-discovery/internal/repo"
	md "github.com/JMURv/service-discovery/pkg/model"
	"go.uber.org/zap"
)

type ServiceDiscoveryRepo interface {
	Register(ctx context.Context, name, addr string) error
	Deregister(ctx context.Context, name, addr string) error
	FindServiceByName(ctx context.Context, name string) (string, error)
	ListServices(ctx context.Context) ([]string, error)
	ListAddrs(ctx context.Context, name string) ([]string, error)
	DeactivateSvc(_ context.Context, name, addr string) error
	ActivateSvc(ctx context.Context, name, addr string) error
	Close() error
}

type Controller struct {
	repo        ServiceDiscoveryRepo
	newAddrChan chan md.Service
}

func New(repo ServiceDiscoveryRepo, newAddrChan chan md.Service) *Controller {
	return &Controller{
		repo:        repo,
		newAddrChan: newAddrChan,
	}
}

func (c *Controller) Register(ctx context.Context, name, addr string) error {
	err := c.repo.Register(ctx, name, addr)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		zap.L().Debug(
			"Error svc already registered",
			zap.String("name", name), zap.String("address", addr),
		)
		return ErrAlreadyExists
	} else if err != nil {
		zap.L().Error(
			"Error registering svc",
			zap.String("name", name), zap.String("address", addr), zap.Error(err),
		)
		return err
	}

	select {
	case c.newAddrChan <- md.Service{Name: name, Address: addr}:
		zap.L().Debug(
			"Sent new address",
			zap.String("name", name), zap.String("address", addr),
		)
	default:
		zap.L().Warn(
			"Channel is full, could not send new address",
			zap.String("name", name), zap.String("address", addr),
		)
	}

	zap.L().Debug(
		"Registered svc",
		zap.String("name", name), zap.String("address", addr),
	)
	return nil
}

func (c *Controller) Deregister(ctx context.Context, name, addr string) error {
	err := c.repo.Deregister(ctx, name, addr)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"Error svc not registered",
			zap.String("name", name), zap.String("address", addr),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Error(
			"Error deregistering svc",
			zap.String("name", name), zap.String("address", addr), zap.Error(err),
		)
		return err
	}

	return nil
}

func (c *Controller) FindServiceByName(ctx context.Context, name string) (string, error) {
	addr, err := c.repo.FindServiceByName(ctx, name)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"Error svc not registered",
			zap.String("name", name),
		)
		return "", ErrNotFound
	} else if err != nil {
		zap.L().Error(
			"Error finding svc",
			zap.String("name", name), zap.Error(err),
		)
		return "", err
	}

	return addr, nil
}

func (c *Controller) ListServices(ctx context.Context) ([]string, error) {
	svcs, err := c.repo.ListServices(ctx)
	if err != nil {
		zap.L().Error("Error finding svcs", zap.Error(err))
		return []string{}, err
	}

	return svcs, nil
}

func (c *Controller) ListAddrs(ctx context.Context, name string) ([]string, error) {
	svcs, err := c.repo.ListAddrs(ctx, name)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug("Error svc not registered")
		return []string{}, ErrNotFound
	} else if err != nil {
		zap.L().Error(
			"Error finding list of addrs",
			zap.Error(err), zap.String("name", name),
		)
		return []string{}, err
	}

	return svcs, nil
}
