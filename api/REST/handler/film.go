package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/pkg"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *Handler) film(w http.ResponseWriter, r *http.Request) {
	var prefix = "/api/film/"

	switch r.Method {
	case "GET":
		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.getFilm(w, id)
	case "PUT":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.putFilm(w, r, id)
	case "PATCH":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.patchFilm(w, r, id)
	case "DELETE":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.deleteFilm(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

}

func (h *Handler) films(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getFilms(w, r.URL.Query().Get("sortBy"))
	case "POST":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		h.createFilm(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) filmSearch(w http.ResponseWriter, r *http.Request) {
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
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(film)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) getFilms(w http.ResponseWriter, sortBy string) {
	films, err := h.services.GetFilms(sortBy)
	if err != nil {
		pkg.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(films)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) createFilm(w http.ResponseWriter, r *http.Request) {
	var request presenter.FilmRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	var id int
	id, err = h.services.CreateFilm(request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%d", id)
}

func (h *Handler) putFilm(w http.ResponseWriter, r *http.Request, id int) {
	var request presenter.FilmRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	actor, err := h.services.PutFilm(id, request)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) patchFilm(w http.ResponseWriter, r *http.Request, id int) {
	var request presenter.FilmRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	if request.Name != nil && len(*request.Name) < 1 || len(*request.Name) > 150 {
		pkg.HandleError(w, errors.New("name length must be in [1; 150]"), http.StatusBadRequest)
		return
	}
	if request.Description != nil && len(*request.Description) > 1000 {
		pkg.HandleError(w, errors.New("description length must be in [0; 1000]"), http.StatusBadRequest)
		return
	}
	if request.Rating != nil && *request.Rating < 0 || *request.Rating > 10 {
		pkg.HandleError(w, errors.New("rating must be in [0; 10]"), http.StatusBadRequest)
		return
	}

	film, err := h.services.PatchFilm(id, request)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(film)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) deleteFilm(w http.ResponseWriter, id int) {
	err := h.services.DeleteFilm(id)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
}

func (h *Handler) searchFilms(w http.ResponseWriter, r *http.Request) {
	var films []presenter.FilmResponse
	var err error
	if r.URL.Query().Get("name") != "" {
		films, err = h.services.SearchFilmsBy("name", r.URL.Query().Get("name"))
	} else if r.URL.Query().Get("actor") != "" {
		films, err = h.services.SearchFilmsBy("actor", r.URL.Query().Get("actor"))
	}
	if err != nil {
		pkg.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(films)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}