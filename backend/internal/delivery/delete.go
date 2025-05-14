package delivery

import (
	"net/http"
	"strconv"
)

func (d *DataDelivery) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue(ID)
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = d.db.Delete(r.Context(), id)
	if err != nil {
		// TODO: если ничего не удалено ответ StatusNoContent
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
