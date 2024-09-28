package main

import (
	"context"
	"fmt"
	"github.com/JMURv/service-discovery/internal/checker"
	"github.com/JMURv/service-discovery/internal/ctrl"
	"github.com/JMURv/service-discovery/internal/hdl"
	"github.com/JMURv/service-discovery/internal/hdl/grpc"
	"github.com/JMURv/service-discovery/internal/hdl/http"
	sqlite "github.com/JMURv/service-discovery/internal/repo/db"
	mem "github.com/JMURv/service-discovery/internal/repo/memory"
	cfg "github.com/JMURv/service-discovery/pkg/config"
	md "github.com/JMURv/service-discovery/pkg/model"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "local.config.yaml"

func mustRegisterLogger(mode string) {
	switch mode {
	case "prod":
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	case "dev":
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Panic("panic occurred", zap.Any("error", err))
			os.Exit(1)
		}
	}()

	conf := cfg.MustLoad(configPath)
	mustRegisterLogger(conf.Server.Mode)

	ctx, cancel := context.WithCancel(context.Background())
	newAddrChan := make(chan md.Service)

	// Setting up main app

	var repo ctrl.ServiceDiscoveryRepo
	switch conf.DB {
	case cfg.InMem:
		repo = mem.New()
	case cfg.SQLite:
		repo = sqlite.New()
	default:
		zap.L().Fatal("Unsupported repo type in configuration")
	}

	check := checker.New(repo, newAddrChan, conf.Checker)
	svc := ctrl.New(repo, newAddrChan)

	var h hdl.Handler
	switch conf.AcceptReq {
	case cfg.HTTP:
		h = http.New(svc)
	case cfg.GRPC:
		h = grpc.New(svc)
	default:
		zap.L().Fatal("Unsupported handler type in configuration")
	}

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c

		zap.L().Info("Shutting down gracefully...")
		cancel()

		repo.Close()
		h.Close()
		os.Exit(0)
	}()

	// Start service
	go check.Start(ctx)
	zap.L().Info(
		fmt.Sprintf("Starting server on %v://%v:%v", conf.Server.Scheme, conf.Server.Domain, conf.Server.Port),
	)
	h.Start(conf.Server.Port)
}
