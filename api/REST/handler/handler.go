package handler

import (
	"filmLibraryVk/internal/service"
	"filmLibraryVk/pkg"
	"github.com/swaggo/http-swagger/v2"
	"net/http"

	"filmLibraryVk/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := initSwagger()

	mux.Handle("/api/actor", pkg.JWTAuthUser(h.actors))
	mux.Handle("/api/actor/", pkg.JWTAuthUser(h.actor))

	mux.Handle("/api/film", pkg.JWTAuthUser(h.films))
	mux.Handle("/api/film/", pkg.JWTAuthUser(h.film))
	mux.Handle("/api/film/search", pkg.JWTAuthUser(h.filmSearch))

	mux.Handle("/api/auth/register", http.HandlerFunc(h.register))
	mux.Handle("/api/auth/authenticate", http.HandlerFunc(h.authenticate))

	mux.Handle("/api/user", pkg.JWTAuthUser(h.users))
	mux.Handle("/api/user/", pkg.JWTAuthUser(h.user))

	//c := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"*"},
	//	AllowCredentials: true,
	//	AllowedMethods:   []string{"*"},
	//	AllowedHeaders:   []string{"*"},
	//	AllowOriginFunc: func(origin string) bool {
	//		return origin == "https://github.com"
	//	},
	//})
	//handler := c.Handler(mux)

	return mux
}

func initSwagger() *http.ServeMux {
	mux := http.NewServeMux()

	docs.SwaggerInfo.Title = "Swagger"
	docs.SwaggerInfo.Description = "This is a distributed calculation server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	return mux
}
