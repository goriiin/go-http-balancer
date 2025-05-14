package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (d *DataDelivery) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue(ID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var req dataDto
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = d.db.Update(r.Context(), id, req.toDomain())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
}
