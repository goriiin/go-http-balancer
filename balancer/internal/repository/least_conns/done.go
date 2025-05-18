package least_conns

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
)

func (p *Pool) Done(be *domain.Backend) {
	if be == nil || be.URL == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	be.DecrementConns()

	if actualBe, ok := p.backendMap[be.URL.String()]; ok && actualBe.Index != -1 && actualBe.Index < p.healthy.Len() {
		p.healthy.Fix(actualBe.Index)
	}
}
