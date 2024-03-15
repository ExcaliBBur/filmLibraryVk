package handler

import (
	"filmLibraryVk/internal/service"
	"net/http"
)

type Handler struct{
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/actor", http.HandlerFunc(h.actor))
	mux.Handle("/api/actor/", http.HandlerFunc(h.actor))

	mux.Handle("/api/film", http.HandlerFunc(h.film))
	mux.Handle("/api/film/", http.HandlerFunc(h.film))

	return mux
}