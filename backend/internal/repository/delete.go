package repository

import (
	"context"
	"fmt"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"log/slog"
)

func (r *DataRepository) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM some_data WHERE id = $1`
	r.log.Info("Deleting data", slog.Int("id", id))

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.log.Error("Failed to delete data",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%w: %v", errs.ErrDeleteFailed, err)
	}

	if result.RowsAffected() == 0 {
		r.log.Error("No data found for deletion ",
			slog.Int("id", id),
		)

		return errs.ErrNotFound
	}

	r.log.Info("Successfully deleted data with ID: %d", id)

	return nil
}
