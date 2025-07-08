package bookings

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
	bookingRouter := router.PathPrefix("/booking").Subrouter()
	bookingRouter.Use(r.authMiddleware.Middleware)

	bookingRouter.HandleFunc("/", r.handler.CreateBooking).Methods(http.MethodPost)
	bookingRouter.HandleFunc("/{id}", r.handler.GetBooking).Methods(http.MethodGet)
}
