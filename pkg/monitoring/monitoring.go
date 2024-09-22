package monitoring

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type monitoring struct {
	reqs         *prometheus.CounterVec
	reqsInFlight prometheus.Gauge
	redisReqs    *prometheus.CounterVec
	latency      *prometheus.HistogramVec

	histograms map[string]*prometheus.HistogramVec
	counters   map[string]*prometheus.CounterVec
	prefix     string
}

type Monitoring interface {
	Register(packageName string)
	Observe(packageName string, method string, value float64)
	Count(packageName string, method string, fail bool)
	Add(packageName string, method string, value int64)

	GetMetricsHandler() http.Handler
	WrapHandler(path string, h http.Handler) http.Handler
}

func New(prefix string) *monitoring {
	m := &monitoring{}
	m.histograms = make(map[string]*prometheus.HistogramVec)
	m.counters = make(map[string]*prometheus.CounterVec)
	m.prefix = prefix

	m.reqs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: prefix + "_requests_total",
			Help: "How many HTTP requests processed",
		},
		[]string{"method", "code"},
	)

	m.reqsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "_pending_requests",
			Help: "How many requests processing at this moment",
		},
	)

	m.latency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    prefix + "_request_duration_milliseconds",
			Help:    "How long it took to process the request",
			Buckets: []float64{0.0001, 0.001, 0.01, 0.1, 0.5, 1, 1.5, 2},
		},
		[]string{"method", "handler"},
	)

	return m
}

func (m *monitoring) Observe(packageName string, method string, value float64) {
	if histogram, ok := m.histograms[packageName]; ok {
		histogram.With(prometheus.Labels{"method": method}).Observe(value)
	}
}

func (m *monitoring) Count(packageName string, method string, fail bool) {
	if counter, ok := m.counters[packageName]; ok {
		if fail {
			method = "fail - " + method
		}
		counter.With(prometheus.Labels{"method": method}).Inc()
	}
}

func (m *monitoring) Add(packageName string, method string, value int64) {
	if counter, ok := m.counters[packageName]; ok {
		counter.With(prometheus.Labels{"method": method}).Add(float64(value))
	}
}

func (m *monitoring) Register(packageName string) {
	kafkaLatency := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_%s_request_duration_milliseconds", m.prefix, packageName),
			Help:    "How long it took to process the stored procedure call",
			Buckets: []float64{0.0001, 0.001, 0.01, 0.1, 0.5, 1, 1.5, 2, 8, 16, 30, 60, 120, 180, 240, 360, 600},
		},
		[]string{"method"},
	)
	kafkaReqs := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_requests_total", m.prefix, packageName),
			Help: "How many Kafka requests processed, partitioned by method/fail",
		},
		[]string{"method"},
	)

	m.histograms[packageName] = kafkaLatency
	m.counters[packageName] = kafkaReqs
}

func (m *monitoring) Reqs() *prometheus.CounterVec {
	return m.reqs
}

func (m *monitoring) ReqsInFlight() prometheus.Gauge {
	return m.reqsInFlight
}

func (m *monitoring) RedisReqs() *prometheus.CounterVec {
	return m.redisReqs
}

func (m *monitoring) Latency() *prometheus.HistogramVec {
	return m.latency
}

func (m *monitoring) GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}

func (m *monitoring) WrapHandler(path string, h http.Handler) http.Handler {
	return promhttp.InstrumentHandlerInFlight(
		m.ReqsInFlight(),
		promhttp.InstrumentHandlerDuration(
			m.Latency().MustCurryWith(prometheus.Labels{"handler": path}),
			promhttp.InstrumentHandlerCounter(m.Reqs(), h),
		),
	)
}
