package httpserver

import (
	"context"
	"fmt"
	"go-clean-template/config"
	"go-clean-template/internal/domain"
	"go-clean-template/internal/facade/httpserver/router"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"
	"net"
	"net/http"
)

type httpServer struct {
	srv           *http.Server
	cancelBaseCtx context.CancelFunc
}

type Provider interface {
	GetService() domain.Service
	GetAppVersion() string
	GetMonitoring() monitoring.Monitoring
	GetLogger() logger.Logger
}

func New(cfg config.HTTP, prov Provider) *httpServer {
	root := router.New(prov)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      root.Router(),
	}

	return &httpServer{
		srv,
		nil,
	}
}

func (h *httpServer) Run(ctx context.Context) error {
	baseCtx, cancel := context.WithCancel(ctx)
	h.cancelBaseCtx = cancel

	h.srv.BaseContext = func(_ net.Listener) context.Context {
		return baseCtx
	}

	return h.srv.ListenAndServe()
}

func (h *httpServer) Stop(ctx context.Context) error {
	h.srv.SetKeepAlivesEnabled(false)
	h.cancelBaseCtx()

	err := h.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("h.srv.Shutdown: %w", err)
	}

	return nil
}

func (h *httpServer) Info() string {
	return h.srv.Addr
}
