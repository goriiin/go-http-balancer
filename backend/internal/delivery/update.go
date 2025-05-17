package delivery

import (
	"encoding/json"
	"errors"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"github.com/goriiin/go-http-balancer/pkg"
	"log/slog"
	"net/http"
	"strconv"
)

func (d *DataDelivery) Update(w http.ResponseWriter, r *http.Request) {
	log := d.log.With(
		slog.String("operation", "[ DataDelivery.Update ]"),
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

	var req dataDto
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)

		return
	}

	log = log.With(slog.Any("request", req))

	if err = d.db.Update(r.Context(), id, req.toDomain()); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			log.Warn("data not found for update")
			pkg.WriteErrorJSON(w, http.StatusNotFound, errs.ErrDataNotFound)

			return
		}

		log.Error("failed to update data",
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusInternalServerError, errs.ErrInternalServer)

		return
	}

	log.Info("data successfully updated")
	pkg.WriteJSON(w, http.StatusNoContent, nil)
}
