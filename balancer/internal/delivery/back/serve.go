package back

import (
	"log/slog"
	"net/http"
)

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.log.Debug("proxying request to backend",
		slog.String("backend_url", b.URL.String()),
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
	)

	b.domainBe.IncrementConns()

	b.Proxy.ServeHTTP(w, r)
}
