package least_conns

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
)

func (p *Pool) Next() *domain.Backend {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.healthy.Len() == 0 {
		return nil
	}

	h := p.healthy.Pop()

	return h
}
