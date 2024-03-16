package repository

import (
	"database/sql"
	"filmLibraryVk/model/dto/actor"
	"filmLibraryVk/model/dto/film"
)

type Actor interface {
	GetActor(id int) (actor.ActorResponse, error)
	GetActors() ([]actor.ActorResponse, error)

	CreateActor(request actor.ActorRequest) (int, error)

	PutActor(id int, request actor.ActorRequest) (actor.ActorResponse, error)
	PatchActor(id int, request actor.ActorRequest) (actor.ActorResponse, error)

	DeleteActor(id int) error
}

type Film interface {
	GetFilm(id int) (film.FilmResponse, error)
	GetFilms(sortBy string) ([]film.FilmResponse, error)

	CreateFilm(request film.FilmRequest) (int, error)

	PutFilm(id int, request film.FilmRequest) (film.FilmResponse, error)
	PatchFilm(id int, request film.FilmRequest) (film.FilmResponse, error)

	DeleteFilm(id int) error

	SearchFilmsByName(name string) ([]film.FilmResponse, error)
	SearchFilmsByActor(name string) ([]film.FilmResponse, error)
}

type Repository struct {
	Actor
	Film
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Actor: NewActorRepo(db),
		Film:  NewFilmRepo(db),
	}
}
