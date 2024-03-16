package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"filmLibraryVk/model/dto/film"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

func (h *Handler) film(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s request on %s", r.Method, r.RequestURI)
	var prefix = "/api/film/"
	switch r.Method {
	case "GET":
		if r.RequestURI != "/api/film/" {
			h.getFilms(w, r.URL.Query().Get("sortBy"))
			return
		}
		id := getPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.getFilm(w, id)
	case "POST":
		h.createFilm(w, r)
	case "PUT":
		id := getPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.putFilm(w, r, id)
	case "PATCH":
		id := getPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.patchFilm(w, r, id)

	case "DELETE":
		id := getPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.deleteFilm(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) filmSearch(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s request on %s", r.Method, r.RequestURI)
	switch r.Method {
	case "GET":
		h.searchFilms(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}


func (h *Handler) getFilm(w http.ResponseWriter, id int) {
	film, err := h.services.GetFilm(id)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(film)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) getFilms(w http.ResponseWriter, sortBy string) {

	films, err := h.services.GetFilms(sortBy)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(films)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) createFilm(w http.ResponseWriter, r *http.Request) {
	var request film.FilmRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	var id int
	id, err = h.services.CreateFilm(request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%d", id)
}

func (h *Handler) putFilm(w http.ResponseWriter, r *http.Request, id int) {
	var request film.FilmRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	actor, err := h.services.PutFilm(id, request)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) patchFilm(w http.ResponseWriter, r *http.Request, id int) {
	var request film.FilmRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if request.Name != nil && len(*request.Name) < 1 || len(*request.Name) > 150 {
		handleError(w, errors.New("name length must be in [1; 150]"), http.StatusBadRequest)
		return
	}
	if request.Description != nil && len(*request.Description) > 1000 {
		handleError(w, errors.New("description length must be in [0; 1000]"), http.StatusBadRequest)
		return
	}
	if request.Rating != nil && *request.Rating < 0 || *request.Rating > 10 {
		handleError(w, errors.New("rating must be in [0; 10]"), http.StatusBadRequest)
		return
	}

	film, err := h.services.PatchFilm(id, request)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(film)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) deleteFilm(w http.ResponseWriter, id int) {
	err := h.services.DeleteFilm(id)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
}

func (h *Handler) searchFilms(w http.ResponseWriter, r *http.Request) {
	var films []film.FilmResponse
	var err error
	if r.URL.Query().Get("name") != "" {
		films, err = h.services.SearchFilmsBy("name", r.URL.Query().Get("name"))
	} else if r.URL.Query().Get("actor") != "" {
		films, err = h.services.SearchFilmsBy("actor", r.URL.Query().Get("actor"))
	}
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(films)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}
