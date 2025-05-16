package least_conns

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"sync/atomic"
)

func (p *Pool) Done(be *domain.Backend) {
	atomic.AddInt32(&be.ConnsCount, -1)
	p.mu.Lock()
	if be.Index >= 0 {
		p.healthy.Fix(be.Index)
	}

	p.mu.Unlock()
}
