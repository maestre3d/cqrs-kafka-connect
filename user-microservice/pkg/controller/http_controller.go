package controller

import "github.com/gorilla/mux"

type HTTPController interface {
	MapRoutes(r *mux.Router)
}
