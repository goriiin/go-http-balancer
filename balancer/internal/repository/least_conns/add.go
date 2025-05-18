package least_conns

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) Add(be *domain.Backend) {
	if be == nil || be.URL == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	urlStr := be.URL.String()
	if _, exists := p.backendMap[urlStr]; exists {
		return
	}

	if _, exists := p.unhealthy[urlStr]; exists {
		return
	}

	be.ResetConns()
	p.healthy.Push(be)
	p.backendMap[urlStr] = be
}
