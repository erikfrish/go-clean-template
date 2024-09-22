package router

import (
	"go-clean-template/internal/domain"
	"go-clean-template/internal/facade/httpserver/handler"
	"go-clean-template/internal/facade/httpserver/middleware"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/monitoring"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

type Provider interface {
	GetService() domain.Service
	GetAppVersion() string
	GetMonitoring() monitoring.Monitoring
	GetLogger() logger.Logger
}

type router struct {
	root *mux.Router
	prov Provider
}

func New(prov Provider) *router {
	root := mux.NewRouter()

	v1Prefix := "/api/v1"
	RegisterDomainHandlers(prov, root, v1Prefix)

	r := router{
		root,
		prov,
	}

	r.initPprofHandlers()
	r.initProbesHandlers()
	r.initMiddlewares()

	return &r
}

func (r *router) Router() http.Handler {
	return r.root
}

func (r *router) initPprofHandlers() {
	debugPrefix := "/debug/pprof"
	r.root.HandleFunc(debugPrefix+"/", pprof.Index)
	r.root.HandleFunc(debugPrefix+"/cmdline", pprof.Cmdline)
	r.root.HandleFunc(debugPrefix+"/symbol", pprof.Symbol)
	r.root.HandleFunc(debugPrefix+"/trace", pprof.Trace)

	profilePrefix := "/profile"
	r.root.HandleFunc(profilePrefix+"", pprof.Profile)
	r.root.Handle(profilePrefix+"/goroutine", pprof.Handler("goroutine"))
	r.root.Handle(profilePrefix+"/threadcreate", pprof.Handler("threadcreate"))
	r.root.Handle(profilePrefix+"/heap", pprof.Handler("heap"))
	r.root.Handle(profilePrefix+"/block", pprof.Handler("block"))
	r.root.Handle(profilePrefix+"/mutex", pprof.Handler("mutex"))
}

func (r *router) initProbesHandlers() {
	h := handler.New(r.prov)

	apiPrefix := "/api"
	r.root.HandleFunc(apiPrefix+"/version", h.GetVersion)
	r.root.HandleFunc(apiPrefix+"/utc", h.GetTimeInUTC)
	r.root.HandleFunc(apiPrefix+"/live", h.GetNoContent)
	r.root.HandleFunc(apiPrefix+"/ready", h.GetNoContent)

	r.root.Handle("/metrics", r.prov.GetMonitoring().GetMetricsHandler())
}

func (r *router) initMiddlewares() {
	mw := middleware.New(r.prov)
	r.root.Use(mw.RecoverMiddleware)
	r.root.Use(mw.RequestLogger)
	r.root.Use(mw.ValidationMiddleware)
	r.root.Use(mw.MonitoringMiddleware)
}
