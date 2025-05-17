package global

import (
	"github.com/goriiin/go-http-balancer/backend/configs"
	"log/slog"
	"net/http"
)

type dataDelivery interface {
	Get(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type App struct {
	dataDelivery dataDelivery
	router       *http.ServeMux
	config       configs.AppConfig
	log          *slog.Logger
	server       *http.Server
}

func NewApp(dataDelivery dataDelivery, router *http.ServeMux, c configs.AppConfig, log *slog.Logger) *App {

	return &App{
		dataDelivery: dataDelivery,
		router:       router,
		config:       c,
		log:          log,
	}
}
