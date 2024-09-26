package main

import (
	"fmt"
	ctrl "github.com/JMURv/service-discovery/internal/ctrl"
	//handler "github.com/JMURv/service-discovery/internal/handler/http"
	handler "github.com/JMURv/service-discovery/internal/hdl/grpc"
	"go.uber.org/zap"
	//mem "github.com/JMURv/service-discovery/internal/repository/memory"
	db "github.com/JMURv/service-discovery/internal/repo/db"
	cfg "github.com/JMURv/service-discovery/pkg/config"
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

	// Setting up main app
	repo := db.New()
	svc := ctrl.New(repo)
	h := handler.New(svc)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c

		zap.L().Info("Shutting down gracefully...")

		repo.Close()
		h.Close()
		os.Exit(0)
	}()

	// Start service
	zap.L().Info(
		fmt.Sprintf("Starting server on %v://%v:%v", conf.Server.Scheme, conf.Server.Domain, conf.Server.Port),
	)
	h.Start(conf.Server.Port)
}
