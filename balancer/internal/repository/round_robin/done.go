package round_robin

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) Done(be *domain.Backend) {
	if be != nil {
		be.DecrementConns()
	}
}
