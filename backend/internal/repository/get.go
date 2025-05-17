package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/goriiin/go-http-balancer/backend/internal/domain"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

func (r *DataRepository) Get(ctx context.Context, id int) (domain.SomeData, error) {
	const query = `SELECT id, name, description FROM some_data WHERE id = $1`
	r.log.Info("Fetching data", slog.Int("id", id))

	var data domain.SomeData
	err := r.db.QueryRow(ctx, query, id).Scan(
		&data.ID,
		&data.Name,
		&data.Description,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Error("Data not found",
				slog.Int("id", id),
				slog.String("error", err.Error()),
			)

			return domain.SomeData{}, errs.ErrNotFound
		}

		r.log.Error("Failed to fetch data",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)

		return domain.SomeData{}, fmt.Errorf("%w: %v", errs.ErrSaveFailed, err)
	}

	r.log.Info("Successfully fetched data", slog.Int("id", id))

	return data, nil
}
