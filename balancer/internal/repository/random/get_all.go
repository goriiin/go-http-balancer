package random

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) GetAll() []*domain.Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	all := make([]*domain.Backend, 0, len(p.healthy)+len(p.unhealthy))
	all = append(all, p.healthy...)

	for _, backend := range p.unhealthy {
		all = append(all, backend)
	}

	return all
}
