package handler

import (
	"encoding/json"
	"errors"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/pkg"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var register presenter.Register
	err := json.NewDecoder(r.Body).Decode(&register)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(register); err != nil {
		pkg.HandleError(w, errors.New("Invalid request body. Username length must be >= 2, " +
			"password length must be [8, 16]"), http.StatusBadRequest)
		return
	}

	jwt, err := h.services.Register(register)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "jwt: %s", jwt)
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var login presenter.Login
	err := json.NewDecoder(r.Body).Decode(&login)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(login); err != nil {
		pkg.HandleError(w, errors.New("Invalid request body. Username length must be >= 2, " +
			"password length must be [8, 16]"), http.StatusBadRequest)
		return
	}

	jwt, err := h.services.Login(login)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "jwt: %s", jwt)
}
