package pkg

import (
	"encoding/json"
	"github.com/goriiin/go-http-balancer/errs"
	"net/http"
)

type Response struct {
	Body interface{} `json:"body"`
}

type HTTPErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		response := Response{
			Body: data,
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, errs.InternalServerError.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func WriteErrorJSON(w http.ResponseWriter, status int, err error) {
	res := HTTPErrorResponse{
		ErrorMessage: err.Error(),
	}

	WriteJSON(w, status, res)
}
