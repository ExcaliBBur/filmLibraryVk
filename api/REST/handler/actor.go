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
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}
		h.getActor(w, id)
	case "PUT":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}
		h.putActor(w, r, id)
	case "PATCH":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}
		h.patchActor(w, r, id)
	case "DELETE":
		if err := pkg.ValidateAdminRoleJWT(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
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

func (h *Handler) mockActors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getActors(w)
	case "POST":
		h.createActor(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) mockActor(w http.ResponseWriter, r *http.Request) {
	var prefix = "/api/actor/"
	switch r.Method {
	case "GET":
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}
		h.getActor(w, id)
	case "PUT":
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}
		h.putActor(w, r, id)
	case "PATCH":
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}
		h.patchActor(w, r, id)
	case "DELETE":
		id, err := pkg.GetPathId(w, r, prefix)
		if err != nil {
			return
		}
		h.deleteActor(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}


// Get actor by id
// @Summary      Get actor by id
// @Description  Get actor by id
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param 		 id   path 	int 	true "id"
// @Success      200  {object}  presenter.ActorResponse
// @Failure      400  {object}  string
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /actor/{id} [get]
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


// Get actors
// @Summary      Get actors
// @Description  Get actors
// @Tags         actors
// @Accept       json
// @Produce      json
// @Success      200  {object}  []presenter.ActorResponse
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /actor [get]
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


// Put actor by id only for ADMIN
// @Summary      Put actor by id
// @Description  Put actor by id
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param 		 id path int true "id"
// @Param 		 request body presenter.ActorRequest true "actor"
// @Success      200  {object}  presenter.ActorResponse
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /actor/{id} [put]
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

// Patch actor by id only for ADMIN
// @Summary      Patch actor by id
// @Description  Patch actor by id
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param 		 id path int true "id"
// @Param 		 request body presenter.ActorRequest true "actor"
// @Success      200  {object}  presenter.ActorResponse
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /actor/{id} [patch]
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

// Delete actor by id only for ADMIN
// @Summary      Delete actor by id
// @Description  Delete actor by id
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param 		 id   path 	int 	true "id"
// @Success      200  {object}  string
// @Failure      401  {object}  string
// @Failure      403  {object}  string
// @Router       /actor/{id} [delete]
func (h *Handler) deleteActor(w http.ResponseWriter, id int) {
	err := h.services.DeleteActor(id)
	if err != nil {
		pkg.HandleError(w, err, http.StatusBadRequest)
		return
	}
}
