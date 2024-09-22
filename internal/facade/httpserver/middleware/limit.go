package middleware

import (
	"net/http"
	"sync"
)

const defaultLimit = 10

var trackingRoutes = map[string]int{ //nolint: gochecknoglobals //tracking routes
	// TODO: move paths to config
	"/api/v1/domain": defaultLimit,
}

type pendingStats struct {
	stats map[string]int
	sync.RWMutex
}

func newPendingStats() *pendingStats {
	return &pendingStats{
		stats: make(map[string]int),
	}
}

func (p *pendingStats) Inc(route string) {
	p.Lock()
	defer p.Unlock()
	p.stats[route]++
}

func (p *pendingStats) Dec(route string) {
	p.Lock()
	defer p.Unlock()
	p.stats[route]--
}

func (p *pendingStats) Get(route string) int {
	p.RLock()
	defer p.RUnlock()
	return p.stats[route]
}

func (m *middleware) LimiterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		m.pendingStats.Inc(path)
		defer m.pendingStats.Dec(path)
		if limit, ok := trackingRoutes[path]; ok {
			if m.pendingStats.Get(path) > limit {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
