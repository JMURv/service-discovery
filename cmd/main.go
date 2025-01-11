package main

import (
	"context"
	"fmt"
	"github.com/JMURv/service-discovery/internal/checker"
	"github.com/JMURv/service-discovery/internal/ctrl"
	"github.com/JMURv/service-discovery/internal/hdl/grpc"
	"github.com/JMURv/service-discovery/internal/hdl/http"
	sqlite "github.com/JMURv/service-discovery/internal/repo/db"
	mem "github.com/JMURv/service-discovery/internal/repo/memory"
	cfg "github.com/JMURv/service-discovery/pkg/config"
	md "github.com/JMURv/service-discovery/pkg/model"
	"go.uber.org/zap"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type Handler interface {
	io.Closer
	Start(port int)
}

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
	defer cancel()

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

	newAddrChan := make(chan md.Service)
	check := checker.New(repo, newAddrChan, conf.Checker)
	svc := ctrl.New(repo, newAddrChan)

	httpH := http.New(svc)
	grpcH := grpc.New(svc)

	// Start service
	zap.L().Info(
		fmt.Sprintf(
			"HTTP is running on %v://%v:%v || GRPC is running on %v://%v:%v",
			conf.Server.Scheme,
			conf.Server.Domain,
			conf.Server.Port,
			conf.Server.Scheme,
			conf.Server.Domain,
			conf.Server.Port+1,
		),
	)
	go check.Start(ctx)
	go httpH.Start(conf.Server.Port)
	go grpcH.Start(conf.Server.Port + 1)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	zap.L().Info("Shutting down gracefully...")

	if err := repo.Close(); err != nil {
		zap.L().Warn("failed to close repo", zap.Error(err))
	}

	if err := httpH.Close(); err != nil {
		zap.L().Warn("failed to close http handler", zap.Error(err))
	}

	if err := grpcH.Close(); err != nil {
		zap.L().Warn("failed to close grpc handler", zap.Error(err))
	}
}
