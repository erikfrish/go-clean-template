package main

import (
	"context"
	"fmt"
	"go-clean-template/config"
	"go-clean-template/internal/app"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load("./config/config.yml")
	if err != nil {
		log.Fatal(fmt.Errorf("config.Load: %w", err))
	}

	logOpts := logger.MakeLoggerOpts(cfg)
	lg := logger.New(logOpts)
	lg.Info(cfg)

	mon := monitoring.New(cfg.PromPrefix)

	application, err := app.New(cfg, mon, lg)
	if err != nil {
		lg.Fatal(fmt.Errorf("app.New: %w", err))
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	errChan := make(chan error)
	go func() {
		errChan <- application.Run(ctx)
	}()

	select {
	case err = <-errChan:
		lg.Error(fmt.Errorf("application.Run: %w", err))
	case exit := <-exitChan:
		lg.Info("SIGNAL:", exit.String())
	}

	timeout := 5 * time.Second //nolint:mnd //5s is enough to shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, timeout)
	defer func() {
		shutdownCancel()
	}()

	err = application.Stop(shutdownCtx)
	if err != nil {
		lg.Error(fmt.Errorf("application.Stop: %w", err))
	}
}
