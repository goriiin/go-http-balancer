package repository

import (
	"context"
	"fmt"
	"github.com/goriiin/go-http-balancer/backend/internal/domain"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"log"
	"log/slog"
)

func (r *DataRepository) Save(ctx context.Context, data domain.SomeData) error {
	const query = `INSERT INTO some_data (name, description) VALUES ($1, $2) RETURNING id`

	r.log.Info("Saving new data")

	err := r.db.QueryRow(
		ctx,
		query,
		data.Name,
		data.Description,
	).Scan(&data.ID)

	if err != nil {
		log.Printf("Failed to save data: %v", err)

		r.log.Error("Failed to delete data",
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%w: %v", errs.ErrSaveFailed, err)
	}

	r.log.Info("Successfully saved data", slog.Int("id", data.ID))

	return nil
}
