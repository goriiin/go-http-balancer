package delivery

import (
	"errors"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"github.com/goriiin/go-http-balancer/pkg"
	"log/slog"
	"net/http"
	"strconv"
)

func (d *DataDelivery) Get(w http.ResponseWriter, r *http.Request) {
	log := d.log.With(
		slog.String("operation", "[ DataDelivery.Get ]"),
	)

	idStr := r.PathValue(ID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid ID format",
			slog.String("id", idStr),
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusBadRequest, errs.ErrInvalidIDParameter)

		return
	}

	log = log.With(slog.Int("id", id))

	data, err := d.db.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			log.Warn("data not found")
			pkg.WriteErrorJSON(w, http.StatusNotFound, errs.ErrDataNotFound)

			return
		}

		log.Error("failed to fetch data",
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusInternalServerError, errs.ErrInternalServer)

		return
	}

	log.Info("data successfully retrieved")
	pkg.WriteJSON(w, http.StatusOK, domainToDto(data))
}
