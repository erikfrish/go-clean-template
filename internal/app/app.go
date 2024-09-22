package app

import (
	"context"
	"fmt"
	"go-clean-template/config"

	"go-clean-template/internal/domain"
	"go-clean-template/internal/facade/cron"
	"go-clean-template/internal/facade/httpserver"
	"go-clean-template/internal/provider"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"

	"golang.org/x/sync/errgroup"
)

type app struct {
	prov       Provider
	httpServer Facade
	cronJob    Facade
	mon        monitoring.Monitoring
	lg         logger.Logger
}

type Provider interface {
	GetService() domain.Service
	GetAppVersion() string
	GetMonitoring() monitoring.Monitoring
	GetLogger() logger.Logger
	Close()
}

type Facade interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	Info() string
}

func New(cfg *config.Config, mon monitoring.Monitoring, lg logger.Logger) (*app, error) {
	prov, err := provider.New(cfg, mon, lg)
	if err != nil {
		return nil, fmt.Errorf("provider.New: %w", err)
	}

	httpServer := httpserver.New(cfg.HTTP, prov)

	cronJob := cron.New(cfg.Schedules, prov)

	return &app{
		prov,
		httpServer,
		cronJob,
		mon,
		lg,
	}, nil
}

func (a *app) Run(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		err := a.httpServer.Run(ctx)
		if err != nil {
			a.lg.Error("httpServer.Run:", err)
		}
		a.lg.Info("http server stopped")
		errChan <- err
	}()
	a.lg.Info("http server started at port", a.httpServer.Info())

	go func() {
		err := a.cronJob.Run(ctx)
		if err != nil {
			a.lg.Error("cronJob.Run:", err)
		}
		a.lg.Info("cron job stopped")
		errChan <- err
	}()
	a.lg.Info("cron job started", a.cronJob.Info())

	return <-errChan
}

func (a *app) Stop(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		err := a.httpServer.Stop(ctx)
		if err != nil {
			return fmt.Errorf("a.httpServer.Stop: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		err := a.cronJob.Stop(ctx)
		if err != nil {
			return fmt.Errorf("a.cronJob.Stop: %w", err)
		}
		return nil
	})

	err := eg.Wait()
	a.lg.Info("Application stopped")
	a.prov.Close()

	return err
}
