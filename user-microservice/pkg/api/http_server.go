package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPApi(r *mux.Router) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
