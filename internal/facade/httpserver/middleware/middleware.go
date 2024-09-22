package middleware

import (
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"
)

type middleware struct {
	pendingStats *pendingStats
	mon          monitoring.Monitoring
	lg           logger.Logger
}

type Provider interface {
	GetMonitoring() monitoring.Monitoring
	GetLogger() logger.Logger
}

func New(prov Provider) *middleware {
	return &middleware{
		newPendingStats(),
		prov.GetMonitoring(),
		prov.GetLogger(),
	}
}
