package user_account

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
	userRouter := router.PathPrefix("/user-account").Subrouter()
	userRouter.Use(r.authMiddleware.Middleware)

	userRouter.HandleFunc("/", r.handler.GetUserAccount).Methods("GET")
}
