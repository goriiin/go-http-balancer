package random

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) Add(be *domain.Backend) {
	if be == nil || be.URL == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	urlStr := be.URL.String()
	if _, exists := p.unhealthy[urlStr]; exists {
		return
	}

	for _, h := range p.healthy {
		if h.URL.String() == urlStr {
			return
		}
	}

	p.healthy = append(p.healthy, be)
}
