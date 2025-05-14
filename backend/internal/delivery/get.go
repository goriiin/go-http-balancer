package delivery

import (
	"github.com/goriiin/go-http-balancer/pkg"
	"net/http"
	"strconv"
)

func (d *DataDelivery) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue(ID)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	data, err := d.db.Get(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	pkg.WriteJSON(w, http.StatusOK, domainToDto(data))
}
