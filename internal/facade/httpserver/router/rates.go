package router

import (
	"go-clean-template/internal/facade/httpserver/handler"

	"github.com/gorilla/mux"
)

func RegisterDomainHandlers(prov Provider, root *mux.Router, prefix string) {
	domainHandler := handler.NewDomainHandler(prov)
	root.HandleFunc("GET "+prefix+"/data", domainHandler.GetObjects)
}
