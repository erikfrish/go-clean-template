package crons

import (
	"fmt"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/schedlock"

	"github.com/robfig/cron/v3"
)

type crons struct {
	c  *cron.Cron
	lg logger.Logger
}

func New(lg logger.Logger) *crons {
	return &crons{
		cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		))),
		lg,
	}
}

func (c *crons) Start() {
	c.c.Run()
}

func (c *crons) Stop() {
	ctx := c.c.Stop()
	<-ctx.Done()
}

func (c *crons) AddCron(spec string, cmd func()) error {
	_, err := c.c.AddFunc(spec, cmd)
	if err != nil {
		return fmt.Errorf("c.AddFunc: %w", err)
	}
	return nil
}

func (c *crons) AddCronWithShedlock(spec string, cmd func(), jobName string, r schedlock.Repository) error {
	_, err := c.c.AddFunc(spec, func() {
		err := schedlock.DoOnce(jobName, cmd, r)
		if err != nil {
			c.lg.Error(err)
		}
	})
	if err != nil {
		return fmt.Errorf("c.AddFunc: %w", err)
	}
	return nil
}
