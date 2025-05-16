package global

import "net/http"

type dataDelivery interface {
	Get(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type App struct {
	dataDelivery dataDelivery
	router       *http.ServeMux
}

func NewApp(dataDelivery dataDelivery, router *http.ServeMux) *App {
	return &App{
		dataDelivery: dataDelivery,
		router:       router,
	}
}
