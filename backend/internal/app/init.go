package app

import "net/http"

type dataDelivery interface {
	Get(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type App struct {
	dataDelivery dataDelivery
}

func NewApp(dataDelivery dataDelivery) *App {
	return &App{
		dataDelivery: dataDelivery,
	}
}
