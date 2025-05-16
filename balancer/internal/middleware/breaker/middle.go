package breaker

import (
	"errors"
	"github.com/goriiin/go-http-balancer/pkg"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (b *Breaker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !b.Allow() {
			pkg.WriteErrorJSON(w, http.StatusServiceUnavailable, errors.New("backend unavailable (circuit-open)"))

			return
		}

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		if rec.status >= 500 {
			b.Failure()
		} else {
			b.Success()
		}
	})
}
