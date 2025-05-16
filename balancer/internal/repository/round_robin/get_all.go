package round_robin

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) GetAll() []*domain.Backend {
	return p.healthy
}
