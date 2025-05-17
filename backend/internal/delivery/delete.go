package delivery

import (
	"errors"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"github.com/goriiin/go-http-balancer/pkg"
	"log/slog"
	"net/http"
	"strconv"
)

func (d *DataDelivery) Delete(w http.ResponseWriter, r *http.Request) {
	log := d.log.With(
		slog.String("operation", "[ DataDelivery.Delete ]"),
	)

	idString := r.PathValue(ID)
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Error("invalid ID format",
			slog.String("id", idString),
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusBadRequest, errs.ErrInvalidIDParameter)

		return
	}

	log = log.With(slog.Int("id", id))

	err = d.db.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			log.Warn("data not found for deletion")
			pkg.WriteErrorJSON(w, http.StatusNotFound, errs.ErrDataNotFound)

			return
		}

		log.Error("failed to delete data",
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusInternalServerError, errs.ErrInternalServer)

		return
	}

	log.Info("data successfully deleted")
	pkg.WriteJSON(w, http.StatusNoContent, nil)
}
