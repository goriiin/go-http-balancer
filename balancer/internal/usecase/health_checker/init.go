package health_checker

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

type pool interface {
	MarkUnhealthy(u *domain.Backend)
	MarkHealthy(u *domain.Backend)
	GetAll() []*domain.Backend
}

type Config struct {
	Path    string //куда тыкаем для чека
	Period  time.Duration
	Timeout time.Duration
}

type HealthChecker struct {
	pool   pool
	config Config
	log    *slog.Logger
	client *http.Client
}

func New(p pool, cfg Config, l *slog.Logger) *HealthChecker {
	return &HealthChecker{
		pool:   p,
		config: cfg,
		log:    l,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}
