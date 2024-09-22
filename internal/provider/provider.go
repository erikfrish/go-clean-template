package provider

import (
	"context"
	"fmt"
	"go-clean-template/config"

	"go-clean-template/internal/domain"
	"go-clean-template/internal/integration/httpclient"
	"go-clean-template/internal/integration/postgres"

	"go-clean-template/internal/service"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type provider struct {
	service domain.Service
	cfg     *config.Config
	mon     monitoring.Monitoring
	lg      logger.Logger
}

func New(cfg *config.Config, mon monitoring.Monitoring, lg logger.Logger) (*provider, error) {
	cl := httpclient.New(cfg.HTTPClient)

	_ = httpclient.NewDataAPI(cl, cfg.API)

	pool, err := postgres.NewPool(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}
	schema := cfg.DB.Schema
	_, err = pool.Exec(context.Background(), fmt.Sprintf("CREATE SCHEMA if not exists %s;", schema))
	if err != nil {
		return nil, fmt.Errorf("failed to create schema %s: %w", schema, err)
	}
	db := stdlib.OpenDBFromPool(pool)
	goose.SetTableName(fmt.Sprintf("%s.goose_db_version", schema))
	err = goose.Up(db, "deploy/migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	lg.Info("connected to database")
	service := service.NewService(lg)

	return &provider{
		service,
		cfg,
		mon,
		lg,
	}, nil
}

func (p *provider) GetService() domain.Service {
	return p.service
}

func (p *provider) GetAppVersion() string {
	return p.cfg.AppVersion
}

func (p *provider) GetMonitoring() monitoring.Monitoring {
	return p.mon
}

func (p *provider) GetLogger() logger.Logger {
	return p.lg
}

func (p *provider) Close() {

}
