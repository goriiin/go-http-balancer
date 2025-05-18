package random

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) MarkHealthy(u *domain.Backend) {
	if u == nil || u.URL == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	urlStr := u.URL.String()
	if backend, exists := p.unhealthy[urlStr]; exists {
		delete(p.unhealthy, urlStr)
		isAlreadyHealthy := false

		for _, h := range p.healthy {
			if h.URL.String() == urlStr {
				isAlreadyHealthy = true

				break
			}
		}

		if !isAlreadyHealthy {
			p.healthy = append(p.healthy, backend)
		}

		return
	}

	for _, h := range p.healthy {
		if h.URL.String() == urlStr {
			return
		}
	}

	p.healthy = append(p.healthy, u)
}
