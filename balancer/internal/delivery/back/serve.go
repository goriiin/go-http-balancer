package back

import (
	"net/http"
)

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.Proxy.ServeHTTP(w, r)
}
