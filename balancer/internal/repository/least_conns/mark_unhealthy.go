package least_conns

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (p *Pool) MarkUnhealthy(u *domain.Backend) {
	if u == nil || u.URL == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	urlStr := u.URL.String()
	if actualBe, ok := p.backendMap[urlStr]; ok && actualBe.Index != -1 {
		removedBe := p.healthy.Remove(actualBe.Index)
		if removedBe != nil {
			delete(p.backendMap, urlStr)
			p.unhealthy[urlStr] = removedBe
		}
	}
}
