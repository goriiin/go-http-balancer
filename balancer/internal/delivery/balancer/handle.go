package balancer

import (
	"errors"
	"github.com/goriiin/go-http-balancer/pkg"
	"log/slog"
	"net/http"
)

func (b *Balancer) Handle(w http.ResponseWriter, r *http.Request) {
	for {
		u := b.pool.Next()
		if u == nil {
			pkg.WriteErrorJSON(w, http.StatusServiceUnavailable, errors.New("no backends"))
			b.log.Error("no backends in health pool", slog.String("place", "[ Balancer.Handle] "))

			return
		}

		b.mux.RLock()
		back, ok := b.m[u.URL]
		b.mux.RUnlock()

		if !ok {
			b.pool.MarkUnhealthy(u)

			continue
		}

		back.ServeHTTP(w, r)

		return
	}
}
