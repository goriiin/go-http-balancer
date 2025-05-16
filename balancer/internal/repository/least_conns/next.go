package least_conns

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"sync/atomic"
)

func (p *Pool) Next() *domain.Backend {
	p.mu.Lock()
	if p.healthy.Len() == 0 {
		p.mu.Unlock()
		return nil
	}

	be := heap.Pop(p.healthy)
	atomic.AddInt32(&be.ConnsCount, 1)
	heap.Push(p.healthy, be)
	p.mu.Unlock()

	return be
}
