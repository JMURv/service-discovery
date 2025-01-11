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

	svcs, err := c.repo.ListServices(ctx)
	if err != nil {
		zap.L().Debug("failed to list services", zap.Error(err))
		return
	}

	for i := 0; i < len(svcs); i++ {
		go c.worker(ctx, svcs[i].Name, svcs[i].Address)
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
			var err error

			switch c.req {
			case config.HTTP:
				err = c.HTTPReq(name, addr)
			case config.GRPC:
				err = c.gRPCReq(name, addr)
			}

			if err != nil {
				zap.L().Debug(
					ErrCheckService.Error(),
					zap.String("svc", name),
					zap.String("addr", addr),
					zap.Error(err),
				)

				if err = c.repo.DeactivateSvc(ctx, name, addr); err != nil {
					zap.L().Debug(
						"failed to deactivate service",
						zap.String("svc", name), zap.String("addr", addr), zap.Error(err),
					)
				}

				c.failedAttempts[name][addr]++
				if c.failedAttempts[name][addr] >= c.conf.MaxRetriesReq {
					zap.L().Info(
						"deregistering due to failed health checks...",
						zap.String("svc", name), zap.String("addr", addr),
					)

					if err = c.repo.Deregister(ctx, name, addr); err != nil {
						zap.L().Debug(
							"failed to deregister service",
							zap.String("svc", name), zap.String("addr", addr), zap.Error(err),
						)
					} else {
						delete(c.failedAttempts[name], addr)
					}
					return
				}
			} else {
				if err = c.repo.ActivateSvc(ctx, name, addr); err != nil {
					zap.L().Debug(
						"failed to activate service",
						zap.String("svc", name), zap.String("addr", addr), zap.Error(err),
					)
				}
				delete(c.failedAttempts[name], addr)
			}

		}
	}
}

func (c *Checker) HTTPReq(name, addr string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/health-check", addr), nil)
	if err != nil {
		zap.L().Debug("failed to create request", zap.Error(err))
		return err
	}

	cli := &http.Client{Timeout: 5 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			zap.L().Error(
				"failed to close response body",
				zap.String("svc", name),
				zap.String("addr", addr),
				zap.Error(err),
			)
		}
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return &ErrUnexpectedStatusCode{resp.StatusCode}
	}
}

func (c *Checker) gRPCReq(name, addr string) error {
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		if err = conn.Close(); err != nil {
			zap.L().Error(
				"failed to close connection",
				zap.String("svc", name),
				zap.String("addr", addr),
				zap.Error(err),
			)
		}
	}()

	check, err := grpc_health_v1.NewHealthClient(conn).
		Check(
			context.Background(), &grpc_health_v1.HealthCheckRequest{
				Service: name,
			},
		)
	if err != nil {
		return err
	}

	if check.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		return ErrNotServing
	} else {
		return nil
	}
}
