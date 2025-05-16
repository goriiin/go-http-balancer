package round_robin

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
)

func (p *Pool) move(u *domain.Backend, src, dst *[]*domain.Backend) {
	if u == nil {
		return
	}

	for i, x := range *src {
		if x.URL.String() == u.URL.String() {
			*dst = append(*dst, u)
			*src = append((*src)[:i], (*src)[i+1:]...)

			return
		}
	}
}
