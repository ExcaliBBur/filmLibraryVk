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

		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}

		h.putUser(w, r, id)
	case "PATCH":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}

		h.patchUser(w, r, id)
	case "DELETE":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}

		h.deleteUser(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) mockUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getUsers(w)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) mockUser(w http.ResponseWriter, r *http.Request) {
	var prefix = "/api/user/"
	switch r.Method {
	case "GET":
		h.getUser(w, r)
	case "PUT":
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}

		h.putUser(w, r, id)
	case "PATCH":
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}

		h.patchUser(w, r, id)
	case "DELETE":
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}

		h.deleteUser(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// Get users
// @Summary      Get users
// @Description  Get users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object} []presenter.UserResponse
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /user [get]
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

// Get user by id
// @Summary      Get user by id
// @Description  Get user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  presenter.UserResponse
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /user/{id} [get]
func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetPathId(w, r, "/api/user/")
	if err != nil {
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

// Put user by id only for ADMIN
// @Summary      Put user by id
// @Description  Put user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 id   path 	int 	true "id"
// @Param 		 request body presenter.UserRequest true "user"
// @Success      200  {object}  presenter.UserResponse
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /user/{id} [put]
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

// Patch user by id only for ADMIN
// @Summary      Patch user by id
// @Description  Patch user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 id   path 	int 	true "id"
// @Param 		 request body presenter.UserRequest true "user"
// @Success      200  {object}  presenter.UserResponse
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /user/{id} [patch]
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

// Delete user by id only for ADMIN
// @Summary      Delete user by id
// @Description  Delete user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 id   path 	int 	true "id"
// @Param 		 request body presenter.UserRequest true "user"
// @Success      200  {object}  string
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /user/{id} [delete]
func (h *Handler) deleteUser(w http.ResponseWriter, id int) {
	err := h.services.DeleteUser(id)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
}
