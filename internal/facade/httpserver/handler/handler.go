package handler

import (
	"encoding/json"
	"go-clean-template/internal/domain"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"
	"net/http"
	"time"
)

type handler struct {
	version string
	start   time.Time
	mon     monitoring.Monitoring
}

type Provider interface {
	GetService() domain.Service
	GetAppVersion() string
	GetMonitoring() monitoring.Monitoring
	GetLogger() logger.Logger
}

func New(prov Provider) *handler {
	return &handler{
		prov.GetAppVersion(),
		time.Now(),
		prov.GetMonitoring(),
	}
}

func (h *handler) GetVersion(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(h.start)

	info := map[string]string{
		"version": h.version,
		"uptime":  uptime.String(),
	}

	err := json.NewEncoder(w).Encode(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) GetTimeInUTC(w http.ResponseWriter, _ *http.Request) {
	utc := map[string]string{
		"utc_dt": time.Now().UTC().Format("2006-01-02T15:04:05"),
	}

	err := json.NewEncoder(w).Encode(&utc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) GetNoContent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
