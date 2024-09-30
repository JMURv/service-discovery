package checker

import (
	"context"
	"fmt"
	"github.com/JMURv/service-discovery/internal/ctrl"
	"github.com/JMURv/service-discovery/pkg/config"
	md "github.com/JMURv/service-discovery/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
	"strings"
	"time"
)

type Checker struct {
	conf           *config.CheckerConfig
	repo           ctrl.ServiceDiscoveryRepo
	newAddrChan    chan md.Service
	failedAttempts map[string]map[string]int
	req            config.AcceptReq
}

func New(repo ctrl.ServiceDiscoveryRepo, newAddr chan md.Service, conf *config.CheckerConfig, req config.AcceptReq) *Checker {
	return &Checker{
		repo:           repo,
		newAddrChan:    newAddr,
		failedAttempts: make(map[string]map[string]int),
		conf:           conf,
		req:            req,
	}
}

func (c *Checker) Start(ctx context.Context) {
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

	zap.L().Info("health check started")
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
			time.Sleep(time.Duration(c.conf.CooldownReq) * time.Second)
			success := false

			switch c.req {
			case config.HTTP:
				success = c.HTTPReq(addr)
			case config.GRPC:
				success = c.gRPCReq(name, addr)
			}

			if !success {
				zap.L().Warn(
					"service health check failed",
					zap.String("svc", name), zap.String("addr", addr),
				)

				if err := c.repo.DeactivateSvc(ctx, name, addr); err != nil {
					zap.L().Error(
						"failed to deactivate service",
						zap.String("svc", name), zap.String("addr", addr), zap.Error(err),
					)
				}

				c.failedAttempts[name][addr]++
				if c.failedAttempts[name][addr] >= c.conf.MaxRetriesReq {
					zap.L().Warn(
						"deregistering service due to failed health checks",
						zap.String("svc", name), zap.String("addr", addr),
					)

					if err := c.repo.Deregister(ctx, name, addr); err != nil {
						zap.L().Error(
							"failed to deregister service",
							zap.String("svc", name), zap.String("addr", addr), zap.Error(err),
						)
					} else {
						delete(c.failedAttempts[name], addr)
					}

					return
				}
			} else {
				if err := c.repo.ActivateSvc(ctx, name, addr); err != nil {
					zap.L().Error(
						"failed to activate service",
						zap.String("svc", name), zap.String("addr", addr), zap.Error(err),
					)
				}
				delete(c.failedAttempts[name], addr)
			}

		}
	}
}

func (c *Checker) HTTPReq(addr string) bool {
	success := false
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/health-check", addr), nil)
	if err != nil {
		zap.L().Debug("failed to create request", zap.Error(err))
		return false
	}

	cli := &http.Client{Timeout: 5 * time.Second}
	resp, err := cli.Do(req)
	if resp != nil {
		if err := resp.Body.Close(); err != nil {
			zap.L().Error("failed to close response body", zap.Error(err))
		}
	}
	if err == nil && resp.StatusCode == http.StatusOK {
		success = true
	}

	return success
}

func (c *Checker) gRPCReq(name, addr string) bool {
	success := false
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Warn("failed to connect to service", zap.String("svc", name), zap.String("addr", addr), zap.Error(err))
		return false
	}

	check, err := grpc_health_v1.NewHealthClient(conn).
		Check(context.Background(), &grpc_health_v1.HealthCheckRequest{
			Service: name,
		})
	if err != nil {
		zap.L().Warn("gRPC health check failed", zap.String("svc", name), zap.String("addr", addr), zap.Error(err))
	} else if check.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
		success = true
	} else {
		zap.L().Warn(
			"service is not in a serving state",
			zap.String("svc", name), zap.String("addr", addr),
			zap.String("status", check.GetStatus().String()),
		)
	}

	if err := conn.Close(); err != nil {
		zap.L().Warn("failed to close connection", zap.String("svc", name), zap.String("addr", addr), zap.Error(err))
	}

	return success
}
