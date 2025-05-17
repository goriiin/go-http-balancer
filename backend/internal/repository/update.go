package repository

import (
	"context"
	"fmt"
	"github.com/goriiin/go-http-balancer/backend/internal/domain"
	"github.com/goriiin/go-http-balancer/backend/internal/errs"
	"log"
)

func (r *DataRepository) Update(ctx context.Context, id int, data domain.SomeData) error {
	const query = `UPDATE some_data SET name = $1, description = $2 WHERE id = $3`
	log.Printf("Updating data with ID: %d, new data: %+v", id, data)

	result, err := r.db.Exec(ctx, query, data.Name, data.Description, id)
	if err != nil {
		log.Printf("Failed to update data (ID: %d): %v", id, err)
		return fmt.Errorf("%w: %v", errs.ErrUpdateFailed, err)
	}

	if result.RowsAffected() == 0 {
		log.Printf("No data found for update (ID: %d)", id)
		return errs.ErrNotFound
	}

	log.Printf("Successfully updated data with ID: %d", id)
	return nil
}
