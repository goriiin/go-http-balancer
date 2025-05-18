package health_checker

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
	"net/http"
	"net/url"
)

func (hc *HealthChecker) checkBackend(be *domain.Backend) {
	probeURL := be.URL.ResolveReference(&url.URL{Path: hc.config.Path})

	req, err := http.NewRequest(http.MethodGet, probeURL.String(), nil)
	if err != nil {
		hc.log.Error("failed to create health check request",
			slog.String("url", probeURL.String()),
			slog.String("error", err.Error()))

		hc.pool.MarkUnhealthy(be)

		return
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		hc.log.Warn("health check probe failed (network error)",
			slog.String("backend", be.URL.String()),
			slog.String("probe_url", probeURL.String()),
			slog.String("error", err.Error()))

		hc.pool.MarkUnhealthy(be)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		hc.log.Debug("health check successful", slog.String("backend", be.URL.String()), slog.Int("status", resp.StatusCode))
		hc.pool.MarkHealthy(be)
	} else {
		hc.log.Warn("health check failed (unhealthy status)", slog.String("backend", be.URL.String()), slog.Int("status", resp.StatusCode))
		hc.pool.MarkUnhealthy(be)
	}
}
