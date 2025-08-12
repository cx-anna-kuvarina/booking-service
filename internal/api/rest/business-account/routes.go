package business_account

import (
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
	bookingRouter := router.PathPrefix("/business-account").Subrouter()
	bookingRouter.Use(r.authMiddleware.Middleware)

	bookingRouter.HandleFunc("/", r.handler.CreateBusinessAccount).Methods("POST")
	bookingRouter.HandleFunc("/{id}", r.handler.UpdateBusinessAccount).Methods("PUT")
	bookingRouter.HandleFunc("/{id}", r.handler.DeleteBusinessAccount).Methods("DELETE")
	bookingRouter.HandleFunc("/{id}", r.handler.GetBusinessAccount).Methods("GET")
}
