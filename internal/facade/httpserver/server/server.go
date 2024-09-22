package server

import (
	"go-clean-template/config"
	"net/http"
)

func New(cfg config.HTTP, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      router,
	}
}
