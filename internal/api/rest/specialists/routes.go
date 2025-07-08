package specialists

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	handler        *Handler
	authMiddleware mux.MiddlewareFunc
}

func NewRouter(handler *Handler, authMiddleware mux.MiddlewareFunc) Router {
	return Router{handler: handler, authMiddleware: authMiddleware}
}

func (r Router) RegisterRoutes(router *mux.Router) {
	specialistRouter := router.PathPrefix("/specialists").Subrouter()
	specialistRouter.Use(r.authMiddleware.Middleware)

	specialistRouter.HandleFunc("/", r.handler.GetSpecialists).Methods(http.MethodGet)
}
