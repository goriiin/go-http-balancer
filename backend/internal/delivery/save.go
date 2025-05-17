package delivery

import (
	"encoding/json"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"github.com/goriiin/go-http-balancer/pkg"
	"log/slog"
	"net/http"
)

func (d *DataDelivery) Save(w http.ResponseWriter, r *http.Request) {
	log := d.log.With(
		slog.String("operation", "[ DataDelivery.Save ]"),
	)

	var req dataDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)

		return
	}

	log = log.With(slog.Any("request", req))

	if err := d.db.Save(r.Context(), req.toDomain()); err != nil {
		log.Error("failed to save data",
			slog.String("error", err.Error()),
		)
		pkg.WriteErrorJSON(w, http.StatusInternalServerError, errs.ErrInternalServer)

		return
	}

	log.Info("data successfully saved")
	pkg.WriteJSON(w, http.StatusCreated, nil)
}
