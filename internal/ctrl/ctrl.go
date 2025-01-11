package ctrl

import (
	"context"
	"errors"
	"github.com/JMURv/service-discovery/internal/repo"
	md "github.com/JMURv/service-discovery/pkg/model"
	"go.uber.org/zap"
	"io"
)

type ServiceDiscoveryRepo interface {
	io.Closer
	ListNames(ctx context.Context) ([]string, error)
	ListAddrsByName(ctx context.Context, name string) ([]string, error)
	ListServices(ctx context.Context) ([]md.Service, error)
	FindServiceByName(ctx context.Context, name string) (string, error)
	Register(ctx context.Context, name, addr string, svcType md.SvcType) error
	Deregister(ctx context.Context, name, addr string) error
	DeactivateSvc(_ context.Context, name, addr string) error
	ActivateSvc(ctx context.Context, name, addr string) error
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

func (c *Controller) ListNames(ctx context.Context) ([]string, error) {
	names, err := c.repo.ListNames(ctx)
	if err != nil {
		zap.L().Error("Error finding svcs", zap.Error(err))
		return nil, err
	}

	return names, nil
}

func (c *Controller) ListAddrsByName(ctx context.Context, name string) ([]string, error) {
	svcs, err := c.repo.ListAddrsByName(ctx, name)
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

func (c *Controller) ListServices(ctx context.Context) ([]md.Service, error) {
	svcs, err := c.repo.ListServices(ctx)
	if err != nil {
		zap.L().Error("Error finding svcs", zap.Error(err))
		return nil, err
	}

	return svcs, nil
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

func (c *Controller) Register(ctx context.Context, name, addr string, svcType md.SvcType) error {
	if err := c.repo.Register(ctx, name, addr, svcType); err != nil && errors.Is(err, repo.ErrAlreadyExists) {
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

	c.newAddrChan <- md.Service{Name: name, Address: addr, SvcType: svcType}
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
