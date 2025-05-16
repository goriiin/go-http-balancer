package least_conns

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) Add(be *domain.Backend) {
	p.mu.Lock()

	heap.Push(p.healthy, be)
	p.mu.Unlock()
}
