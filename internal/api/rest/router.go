package rest

import (
	"github.com/gorilla/mux"
)

type Register interface {
	RegisterRoutes(route *mux.Router)
}

func NewRouter(routes []Register) *mux.Router {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	for _, r := range routes {
		r.RegisterRoutes(apiRouter)
	}

	return router
}
