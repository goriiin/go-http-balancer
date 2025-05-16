package global

import (
	"errors"
	"github.com/goriiin/go-http-balancer/pkg"
	"net/http"
)

func (a *App) crud(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.dataDelivery.Get(w, r)
	case http.MethodPost:
		a.dataDelivery.Save(w, r)
	case http.MethodPut:
		a.dataDelivery.Update(w, r)
	case http.MethodDelete:
		a.dataDelivery.Delete(w, r)
	default:
		pkg.WriteErrorJSON(w, http.StatusMethodNotAllowed, errors.New("status not allowed"))
	}
}
