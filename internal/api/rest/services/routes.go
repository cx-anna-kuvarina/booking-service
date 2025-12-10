package services

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
	servicesRouter := router.PathPrefix("/services").Subrouter()
	servicesRouter.Use(r.authMiddleware.Middleware)

	// CRUD operations for services
	servicesRouter.HandleFunc("/", r.handler.CreateService).Methods(http.MethodPost)
	servicesRouter.HandleFunc("/", r.handler.ListServices).Methods(http.MethodGet)
	servicesRouter.HandleFunc("/{id}", r.handler.GetService).Methods(http.MethodGet)
	servicesRouter.HandleFunc("/{id}", r.handler.UpdateService).Methods(http.MethodPut)
	servicesRouter.HandleFunc("/{id}", r.handler.DeleteService).Methods(http.MethodDelete)

	// Get services by business account
	servicesRouter.HandleFunc("/business-account/{business_account_id}", r.handler.GetServicesByBusinessAccount).Methods(http.MethodGet)
}
