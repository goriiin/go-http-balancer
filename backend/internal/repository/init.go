package repository

import (
	"github.com/goriiin/go-http-balancer/pkg/postgtresql"
	"log/slog"
)

type DataRepository struct {
	db  postgtresql.DB
	log *slog.Logger
}

func NewDataRepository(db postgtresql.DB, log *slog.Logger) *DataRepository {
	return &DataRepository{
		db:  db,
		log: log,
	}
}
