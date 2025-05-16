package round_robin

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
)

func (p *Pool) MarkUnhealthy(u *domain.Backend) {
	if u == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.move(u, &p.healthy, &p.unhealthy)
}
