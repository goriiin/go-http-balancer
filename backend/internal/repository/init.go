package repository

import "github.com/goriiin/go-http-balancer/backend/db/postgtresql"

type DataRepository struct {
	db postgtresql.DB
}
