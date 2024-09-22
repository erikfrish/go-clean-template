package httpclient

import (
	"net/http"

	"go-clean-template/config"
)

func New(cfg config.HTTPClient) *http.Client {
	return &http.Client{
		Timeout: cfg.Timeout,
	}
}
