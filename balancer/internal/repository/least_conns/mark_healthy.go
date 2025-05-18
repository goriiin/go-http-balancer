package least_conns

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) MarkHealthy(u *domain.Backend) {
	if u == nil || u.URL == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	urlStr := u.URL.String()
	if backend, exists := p.unhealthy[urlStr]; exists {
		delete(p.unhealthy, urlStr)

		backend.ResetConns()
		p.healthy.Push(backend)
		p.backendMap[urlStr] = backend

		return
	}

	if _, exists := p.backendMap[urlStr]; exists {
		return
	}

	u.ResetConns()
	p.healthy.Push(u)
	p.backendMap[urlStr] = u
}
