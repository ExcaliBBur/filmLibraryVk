package handler

import (
	"filmLibraryVk/internal/service"
	"filmLibraryVk/pkg"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/api/actor", pkg.JWTAuthUser(h.actors))
	mux.Handle("/api/actor/", pkg.JWTAuthUser(h.actor))

	mux.Handle("/api/film", pkg.JWTAuthUser(h.films))
	mux.Handle("/api/film/", pkg.JWTAuthUser(h.film))
	mux.Handle("/api/film/search", pkg.JWTAuthUser(h.filmSearch))

	mux.Handle("/api/auth/register", http.HandlerFunc(h.register))
	mux.Handle("/api/auth/authenticate", http.HandlerFunc(h.authenticate))

	mux.Handle("/api/user", pkg.JWTAuthUser(h.users))
	mux.Handle("/api/user/", pkg.JWTAuthUser(h.user))

	return mux
}
