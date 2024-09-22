package cron

import (
	"context"
	"fmt"
	"go-clean-template/config"
	"go-clean-template/internal/domain"
	"go-clean-template/pkg/crons"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"
	"time"
)

type cron struct {
	scheds        config.Schedules
	service       domain.Service
	cs            Crons
	baseCtx       context.Context
	cancelBaseCtx context.CancelFunc
	lg            logger.Logger
}

type Provider interface {
	GetService() domain.Service
	GetAppVersion() string
	GetMonitoring() monitoring.Monitoring
	GetLogger() logger.Logger
}

type Crons interface {
	Start()
	Stop()
	AddCron(spec string, cmd func()) error
}

func New(cfg config.Schedules, prov Provider) *cron {
	return &cron{
		cfg,
		prov.GetService(),
		crons.New(prov.GetLogger()),
		nil,
		nil,
		prov.GetLogger(),
	}
}

func (c *cron) persistYesterdayData() {
	const op = "cron.persistYesterdayData"
	tn := time.Now()
	c.lg.Info(fmt.Sprintf("%s: %v", op, tn))

	yesterday := tn.AddDate(0, 0, -1).Format(time.DateOnly)

	err := c.service.Persist(c.baseCtx, yesterday)
	if err != nil {
		c.lg.Error(fmt.Errorf("%s: %w", op, err))
	}

	c.lg.Info(fmt.Sprintf("%s done in %v", op, time.Since(tn)))
}

func (c *cron) Run(ctx context.Context) error {
	c.baseCtx, c.cancelBaseCtx = context.WithCancel(ctx)

	err := c.cs.AddCron(c.scheds.Persist, c.persistYesterdayData)
	if err != nil {
		c.lg.Error("failed to add cron", err)
	}

	c.cs.Start()
	return nil
}

func (c *cron) Stop(ctx context.Context) error {
	c.cancelBaseCtx()
	c.cs.Stop()
	return nil
}

func (c *cron) Info() string {
	return "cron"
}
