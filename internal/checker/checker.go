package checker

import (
	"context"
	"fmt"
	"github.com/JMURv/service-discovery/internal/ctrl"
	md "github.com/JMURv/service-discovery/pkg/model"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// TODO: To config
const maxRetriesReq = 3
const cooldownReq = 5 * time.Second

type Checker struct {
	repo           ctrl.ServiceDiscoveryRepo
	newAddrChan    chan md.Service
	failedAttempts map[string]map[string]int
}

func New(repo ctrl.ServiceDiscoveryRepo, newAddr chan md.Service) *Checker {
	return &Checker{
		repo:           repo,
		newAddrChan:    newAddr,
		failedAttempts: make(map[string]map[string]int),
	}
}

func (c *Checker) healthCheck(ctx context.Context) {
	go c.listenForNewAddresses(ctx)

	names, err := c.repo.ListServices(ctx)
	if err != nil {
		zap.L().Debug("failed to list services", zap.Error(err))
		return
	}

	for _, name := range names {
		addrs, err := c.repo.ListAddrs(ctx, name)
		if err != nil {
			zap.L().Debug("failed to list addrs", zap.Error(err))
			continue
		}

		for _, addr := range addrs {
			go c.worker(ctx, name, addr)
		}
	}

	select {
	case <-ctx.Done():
		zap.L().Info("health check stopped")
		return
	}
}

func (c *Checker) listenForNewAddresses(ctx context.Context) {
	for newSvc := range c.newAddrChan {
		go c.worker(ctx, newSvc.Name, newSvc.Address)
	}
}

func (c *Checker) worker(ctx context.Context, name, addr string) {
	if _, exists := c.failedAttempts[name]; !exists {
		c.failedAttempts[name] = make(map[string]int)
	}

	for {
		select {
		case <-ctx.Done():
			zap.L().Info("worker stopped", zap.String("svc", name), zap.String("addr", addr))
			return
		default:
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/health-check", addr), nil)
			if err != nil {
				zap.L().Debug("failed to create request", zap.Error(err))
				return
			}

			cli := &http.Client{Timeout: 5 * time.Second}
			resp, err := cli.Do(req)
			if err != nil || resp.StatusCode >= 300 {
				zap.L().Warn(
					"service health check failed",
					zap.String("svc", name), zap.String("addr", addr), zap.Int("status", resp.StatusCode),
				)

				c.failedAttempts[name][addr]++
				if c.failedAttempts[name][addr] >= maxRetriesReq {
					zap.L().Warn("deregistering service due to failed health checks", zap.String("svc", name), zap.String("addr", addr))
					if err := c.repo.Deregister(ctx, name, addr); err != nil {
						zap.L().Error("failed to deregister service", zap.String("svc", name), zap.String("addr", addr), zap.Error(err))
					} else {
						delete(c.failedAttempts[name], addr)
					}
					break
				}
			} else {
				delete(c.failedAttempts[name], addr)
			}

			if resp != nil {
				if err := resp.Body.Close(); err != nil {
					zap.L().Error("failed to close response body", zap.Error(err))
				}
			}

			time.Sleep(cooldownReq)
		}
	}
}
