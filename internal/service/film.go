package service

import (
	"filmLibraryVk/internal/repository"
	"filmLibraryVk/model/dto/film"
)

type FilmService struct {
	repo repository.Film
}

func NewFilmService(repo repository.Film) *FilmService {
	return &FilmService{
		repo: repo,
	}
}

func (s *FilmService) GetFilm(id int) (film.FilmResponse, error) {
	return s.repo.GetFilm(id)
}
func (s *FilmService) GetFilms() ([]film.FilmResponse, error) {
	return s.repo.GetFilms()
}

func (s *FilmService) CreateFilm(request film.FilmRequest) (int, error) {
	return s.repo.CreateFilm(request)
}

func (s *FilmService) PutFilm(id int, request film.FilmRequest) (film.FilmResponse, error) {
	return s.repo.PutFilm(id, request)
}

func (s *FilmService) PatchFilm(id int, request film.FilmRequest) (film.FilmResponse, error) {
	return s.repo.PatchFilm(id, request)
}

func (s *FilmService) DeleteFilm(id int) error {
	return s.repo.DeleteFilm(id)
}