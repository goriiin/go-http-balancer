package delivery

import (
	"encoding/json"
	"net/http"
)

func (d *DataDelivery) Save(w http.ResponseWriter, r *http.Request) {
	var req dataDto

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = d.db.Save(r.Context(), req.toDomain())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
