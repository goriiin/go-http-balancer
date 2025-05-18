package health_checker

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
	"sync"
)

func (hc *HealthChecker) performChecks() {
	backends := hc.pool.GetAll()
	if len(backends) == 0 {
		hc.log.Debug("no backends to health check")

		return
	}

	hc.log.Debug("performing health checks",
		slog.Int("backend_count", len(backends)))

	var wg sync.WaitGroup
	for _, be := range backends {
		if be == nil || be.URL == nil {
			hc.log.Warn("nil backend or backend URL found in pool during health check")

			continue
		}

		wg.Add(1)
		go func(backendToCheck *domain.Backend) {
			defer wg.Done()

			hc.checkBackend(backendToCheck)
		}(be)
	}

	wg.Wait()
}
