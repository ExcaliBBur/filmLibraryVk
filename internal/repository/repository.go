package repository

import (
	"database/sql"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/internal/model/entity"
)

type Actor interface {
	GetActor(id int) (presenter.ActorResponse, error)
	GetActors() ([]presenter.ActorResponse, error)

	CreateActor(request presenter.ActorRequest) (int, error)

	PutActor(id int, request presenter.ActorRequest) (presenter.ActorResponse, error)
	PatchActor(id int, request presenter.ActorRequest) (presenter.ActorResponse, error)

	DeleteActor(id int) error
}

type Film interface {
	GetFilm(id int) (presenter.FilmResponse, error)
	GetFilms(sortBy string) ([]presenter.FilmResponse, error)

	CreateFilm(request presenter.FilmRequest) (int, error)

	PutFilm(id int, request presenter.FilmRequest) (presenter.FilmResponse, error)
	PatchFilm(id int, request presenter.FilmRequest) (presenter.FilmResponse, error)

	DeleteFilm(id int) error

	SearchFilmsByName(name string) ([]presenter.FilmResponse, error)
	SearchFilmsByActor(name string) ([]presenter.FilmResponse, error)
}

type User interface {
	GetUserById(id int) (presenter.UserResponse, error)
	GetUserByUsername(username string) (entity.User, error)
	GetUsers() ([]presenter.UserResponse, error)

	PutUser(id int, request presenter.UserRequest) (presenter.UserResponse, error)
	PatchUser(id int, request presenter.UserRequest) (presenter.UserResponse, error)

	DeleteUser(id int) error

	CreateUser(register presenter.Register) (int, error)
}

type Repository struct {
	Actor
	Film
	User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Actor: NewActorRepo(db),
		Film:  NewFilmRepo(db),
		User:  NewUserRepo(db),
	}
}
