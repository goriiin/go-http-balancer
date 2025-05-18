package balancer

import (
	"errors"
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
	"net/http"

	"github.com/goriiin/go-http-balancer/pkg" // For WriteErrorJSON
)

func (b *Balancer) Handle(w http.ResponseWriter, r *http.Request) {
	var selectedHandler http.Handler
	var chosenDomainBe *domain.Backend

	b.mu.RLock()
	numConfiguredHandlers := len(b.backends)
	b.mu.RUnlock()

	if numConfiguredHandlers == 0 {
		b.log.Error("no backends configured to serve request",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		)
		pkg.WriteErrorJSON(w, http.StatusServiceUnavailable, errors.New("no backends available"))
		return
	}

	for attempts := 0; attempts < numConfiguredHandlers; attempts++ {
		chosenDomainBe = b.pool.Next()
		if chosenDomainBe == nil {
			b.log.Warn("balancer pool returned no backend candidate on this attempt",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.Int("attempt", attempts+1),
			)

			continue
		}

		backendURLStr := chosenDomainBe.URL.String()
		b.mu.RLock()
		handler, ok := b.backends[backendURLStr]
		b.mu.RUnlock()

		if !ok {
			b.log.Error("consistency error: backend from pool not found in balancer's handler map. Marking unhealthy.",
				slog.String("backend_url_from_pool", backendURLStr),
				slog.String("path", r.URL.Path),
			)

			b.pool.MarkUnhealthy(chosenDomainBe)

			chosenDomainBe.ResetConns()
			chosenDomainBe = nil

			continue
		}

		selectedHandler = handler

		break
	}

	if selectedHandler == nil {
		b.log.Error("failed to route request: no healthy and mapped backend found after attempts",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.Int("max_attempts_made", numConfiguredHandlers),
		)

		pkg.WriteErrorJSON(w, http.StatusServiceUnavailable, errors.New("all backends are currently unavailable or misconfigured"))

		return
	}

	b.log.Debug("routing request to backend",
		slog.String("backend_url", chosenDomainBe.URL.String()),
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
	)

	selectedHandler.ServeHTTP(w, r)

	b.pool.Done(chosenDomainBe)
}
