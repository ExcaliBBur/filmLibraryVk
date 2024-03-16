package service

import (
	"errors"
	"filmLibraryVk/internal/repository"
	"filmLibraryVk/model/dto/film"
	"filmLibraryVk/model/entity"
	"fmt"
	"reflect"
	"strings"
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
func (s *FilmService) GetFilms(sortBy string) ([]film.FilmResponse, error) {
	sortQuery, err := validateAndReturnSortQuery(sortBy)
	if err != nil {
		return nil, err
	}
	return s.repo.GetFilms(sortQuery)
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

func (s *FilmService) SearchFilmsBy(field, value string) ([]film.FilmResponse, error) {
	switch field{
	case "name":
		return s.repo.SearchFilmsByName(value)
	case "actor" :
		return s.repo.SearchFilmsByActor(value)
	default:
		return nil, errors.New("can not search by " + field)
	}
}

var filmFields = getFilmFields()

func getFilmFields() []string {
	var field []string
	v := reflect.ValueOf(entity.Film{})
	for i := 0; i < v.Type().NumField(); i++ {
		field = append(field, v.Type().Field(i).Tag.Get("json"))
	}
	return field
}

func stringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}

func validateAndReturnSortQuery(sortBy string) (string, error) {
	if sortBy == "" {
		sortBy = "rating.desc"
	}

	splits := strings.Split(sortBy, ".")
	if len(splits) != 2 {
		return "", errors.New("malformed sortBy query parameter, should be field.orderdirection")
	}
	field, order := splits[0], splits[1]
	if order != "desc" && order != "asc" {
		return "", errors.New("malformed orderdirection in sortBy query parameter, should be asc or desc")
	}
	if !stringInSlice(filmFields, field) {
		return "", errors.New("unknown field in sortBy query parameter")
	}

	if strings.ToLower(field) == "releasedate" {
		field = "release_date"
	}
	return fmt.Sprintf("%s %s", field, strings.ToUpper(order)), nil
}
