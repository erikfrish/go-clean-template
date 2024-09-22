package handler

import (
	"errors"
	"go-clean-template/internal/domain"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"

	"encoding/json"

	"net/http"
)

type domainHandler struct {
	service domain.Service
	mon     monitoring.Monitoring
	lg      logger.Logger
}

func NewDomainHandler(prov Provider) *domainHandler {
	return &domainHandler{
		prov.GetService(),
		prov.GetMonitoring(),
		prov.GetLogger(),
	}
}

func (h *domainHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	req := domain.ServiceRequest{}

	err := req.Validate()
	if err != nil {
		h.lg.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Do(r.Context(), req)
	if err != nil {
		if errors.As(err, &domain.ErrNotFound) {
			h.lg.Error("Error getting objects", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		h.lg.Error("Error getting objects", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Data string `json:"data"`
	}{}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		h.lg.Error("Error encoding response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
