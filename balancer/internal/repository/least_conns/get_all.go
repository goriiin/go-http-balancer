package least_conns

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) GetAll() []*domain.Backend {
	p.mu.Lock()
	defer p.mu.Unlock()

	all := make([]*domain.Backend, 0, p.healthy.Len()+len(p.unhealthy))

	healthyBackends := p.healthy.GetAll()
	all = append(all, healthyBackends...)

	for _, backend := range p.unhealthy {
		all = append(all, backend)
	}

	return all
}
