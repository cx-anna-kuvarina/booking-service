package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) Router {
	return Router{handler: handler}
}

func (r Router) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/google-login", r.handler.GoogleLogin).Methods(http.MethodPost)
	router.HandleFunc("/google-callback", r.handler.GoogleCallback).Methods(http.MethodPost)
}
