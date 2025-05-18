package random

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) MarkUnhealthy(u *domain.Backend) {
	if u == nil || u.URL == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	urlStr := u.URL.String()
	newHealthy := make([]*domain.Backend, 0, len(p.healthy))
	found := false

	for _, backend := range p.healthy {
		if backend.URL.String() == urlStr {
			found = true
		} else {
			newHealthy = append(newHealthy, backend)
		}
	}

	if found {
		p.healthy = newHealthy
		p.unhealthy[urlStr] = u
	}
}
