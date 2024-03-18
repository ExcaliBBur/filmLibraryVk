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

func (h *Handler) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getUsers(w)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) user(w http.ResponseWriter, r *http.Request) {
	var prefix = "/api/user/"
	switch r.Method {
	case "GET":
		h.getUser(w, r)
	case "PUT":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}

		h.putUser(w, r, id)
	case "PATCH":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}

		h.patchUser(w, r, id)
	case "DELETE":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}

		h.deleteUser(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getUsers(w http.ResponseWriter) {
	users, err := h.services.GetUsers()
	if err != nil {
		pkg.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(users)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id := pkg.GetPathId(w, r, "/api/user/")
	if id == -1 {
		return
	}

	user, err := h.services.GetUserById(id)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(user)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) putUser(w http.ResponseWriter, r *http.Request, id int) {
	var request presenter.UserRequest
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

	actor, err := h.services.PutUser(id, request)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) patchUser(w http.ResponseWriter, r *http.Request, id int) {
	var request presenter.UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	if request.Username != nil && len(*request.Username) < 2 {
		pkg.HandleError(w, errors.New("username length must be greater than 2"), http.StatusBadRequest)
		return
	}
	if request.Password != nil && len(*request.Password) < 8 {
		pkg.HandleError(w, errors.New("password length must be greater than 8"), http.StatusBadRequest)
		return
	}

	film, err := h.services.PatchUser(id, request)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(film)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) deleteUser(w http.ResponseWriter, id int) {
	err := h.services.DeleteUser(id)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
}