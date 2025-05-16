package health_checker

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"time"
)

type pool interface {
	MarkUnhealthy(u *domain.Backend)
	MarkHealthy(u *domain.Backend)
	GetAll() []*domain.Backend
}

type HealthChecker struct {
	Pool    pool
	Path    string
	Period  time.Duration
	Timeout time.Duration
}

func NewHealthChecker(pool pool) *HealthChecker {
	return &HealthChecker{
		Pool: pool,
	}
}
