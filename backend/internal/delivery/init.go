package delivery

import (
	"context"
	"github.com/goriiin/go-http-balancer/backend/internal/domain"
)

const ID = "id"

type db interface {
	Get(ctx context.Context, id int) (domain.SomeData, error)
	Save(ctx context.Context, data domain.SomeData) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, data domain.SomeData) error
}

// TODO: возвращать более корректные ошибки

type DataDelivery struct {
	db db
}
