package random

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) Next() *domain.Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.healthy) == 0 {
		return nil
	}

	if len(p.healthy) == 1 {
		return p.healthy[0]
	}

	randomIndex := p.rng.Intn(len(p.healthy))

	return p.healthy[randomIndex]
}
