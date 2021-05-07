package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/p44elovod/auth-with-gopg/pkg/service"
)

type Handler struct {
	ctx      context.Context
	services *service.Service
}

func NewHandler(ctx context.Context, services *service.Service) *Handler {
	return &Handler{
		services: services,
		ctx:      ctx,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	auth := router.Methods("POST").Subrouter()
	auth.HandleFunc("/sign-up", h.signUp)
	auth.HandleFunc("/sign-in", h.signIn)

	api := router.Methods(http.MethodGet).Subrouter()
	api.HandleFunc("/api", h.helloAPI)
	api.Use(h.authMiddleWare)

	router.HandleFunc("/refresh", h.refreshTokenPair).Methods("POST")

	return router
}
