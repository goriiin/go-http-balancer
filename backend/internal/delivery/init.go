package delivery

import (
	"context"
	"github.com/goriiin/go-http-balancer/backend/internal/domain"
	"log/slog"
)

const ID = "id"

type db interface {
	Get(ctx context.Context, id int) (domain.SomeData, error)
	Save(ctx context.Context, data domain.SomeData) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, data domain.SomeData) error
}

type DataDelivery struct {
	log *slog.Logger
	db  db
}

func NewDataDelivery(db db, log *slog.Logger) *DataDelivery {
	return &DataDelivery{
		db:  db,
		log: log,
	}
}
