package handler

import (
	"bytes"
	"encoding/json"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/pkg"
	"fmt"
	"net/http"
)

func (h *Handler) actor(w http.ResponseWriter, r *http.Request) {
	var prefix = "/api/actor/"
	switch r.Method {
	case "GET":
		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.getActor(w, id)
	case "PUT":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.putActor(w, r, id)
	case "PATCH":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.patchActor(w, r, id)
	case "DELETE":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id := pkg.GetPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.deleteActor(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) actors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getActors(w)
	case "POST":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		h.createActor(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getActor(w http.ResponseWriter, id int) {
	actor, err := h.services.GetActor(id)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) getActors(w http.ResponseWriter) {
	actors, err := h.services.GetActors()
	if err != nil {
		pkg.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actors)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) createActor(w http.ResponseWriter, r *http.Request) {
	var request presenter.ActorRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	var id int
	id, err = h.services.CreateActor(request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%d", id)
}

func (h *Handler) putActor(w http.ResponseWriter, r *http.Request, id int) {
	var request presenter.ActorRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	actor, err := h.services.PutActor(id, request)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) patchActor(w http.ResponseWriter, r *http.Request, id int) {
	var request presenter.ActorRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	actor, err := h.services.PatchActor(id, request)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) deleteActor(w http.ResponseWriter, id int) {
	err := h.services.DeleteActor(id)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
}
