package service

import (
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/internal/repository"
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

	SearchFilmsBy(field, value string) ([]presenter.FilmResponse, error)
}

type User interface {
	GetUserById(id int) (presenter.UserResponse, error)
	GetUsers() ([]presenter.UserResponse, error)

	PutUser(id int, request presenter.UserRequest) (presenter.UserResponse, error)
	PatchUser(id int, request presenter.UserRequest) (presenter.UserResponse, error)

	DeleteUser(id int) error

	Login(login presenter.Login) (string, error)

	Register(register presenter.Register) (string, error)
}

type Service struct {
	Actor
	Film
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Actor: NewActorService(repo.Actor),
		Film:  NewFilmService(repo.Film),
		User:  NewUserService(repo.User),
	}
}
