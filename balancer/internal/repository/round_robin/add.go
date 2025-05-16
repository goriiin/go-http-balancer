package round_robin

import (
	"errors"
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
)

func (p *Pool) Add(u *domain.Backend) error {
	if u == nil {
		return errors.New("nil backend")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.healthy = append(p.healthy, u)

	return nil
}
