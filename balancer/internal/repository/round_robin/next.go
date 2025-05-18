package round_robin

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
)

func (p *Pool) Next() *domain.Backend {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.healthy) == 0 {
		return nil
	}

	be := p.healthy[p.pos%uint32(len(p.healthy))]
	p.pos++

	return be
}
