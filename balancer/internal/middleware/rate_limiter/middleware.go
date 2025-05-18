package rate_limiter

import (
	"errors"
	"fmt"
	"github.com/goriiin/go-http-balancer/balancer/configs"
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/token_bucket"
	"github.com/goriiin/go-http-balancer/pkg"
	"log/slog"
	"net/http"
)

type Store interface {
	GetBucket(clientID string) (*token_bucket.TokenBucket, error)
	AddClientConfig(config domain.ClientInfo) error
	UpdateClientConfig(config domain.ClientInfo) error
	DeleteClientConfig(clientID string) error
	Stop()
}

func Middleware(store Store, config configs.RateLimiterConfig, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := getClientID(r, config, logger)

			bucket, err := store.GetBucket(clientID)
			if err != nil {
				logger.Error("failed to get token bucket for client",
					slog.String("client_id", clientID),
					slog.String("error", err.Error()),
				)

				pkg.WriteErrorJSON(w, http.StatusInternalServerError, errors.New("rate limiter internal error"))

				return
			}

			if !bucket.Allow() {
				logger.Warn("rate limit exceeded for client",
					slog.String("client_id", clientID),
					slog.String("path", r.URL.Path),
					slog.Int64("bucket_tokens", bucket.CurrentTokens()),
					slog.Int64("bucket_capacity", bucket.Capacity()),
				)

				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", bucket.Capacity()))
				w.Header().Set("X-RateLimit-Remaining", "0")

				pkg.WriteErrorJSON(w, http.StatusTooManyRequests,
					errors.New("rate limit exceeded"))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
