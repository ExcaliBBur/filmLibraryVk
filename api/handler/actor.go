package handler

import (
	"bytes"
	"encoding/json"
	"filmLibraryVk/model/dto/actor"
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) actor(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s request on %s", r.Method, r.RequestURI)
	var prefix = "/api/actor/"
	switch r.Method {
	case "GET":
		if r.RequestURI != "/api/actor" {
			id := getPathId(w, r, prefix)
			if id == -1 {
				return
			}
			h.getActor(w, id)
			return
		}

		h.getActors(w)
	case "POST":
		h.createActor(w, r)
	case "PUT":
		id := getPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.putActor(w, r, id)
	case "PATCH":
		id := getPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.patchActor(w, r, id)

	case "DELETE":
		id := getPathId(w, r, prefix)
		if id == -1 {
			return
		}
		h.deleteActor(w, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getActor(w http.ResponseWriter, id int) {
	actor, err := h.services.GetActor(id)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) getActors(w http.ResponseWriter) {
	actors, err := h.services.GetActors()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actors)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) createActor(w http.ResponseWriter, r *http.Request) {
	var request actor.ActorRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	var id int
	id, err = h.services.CreateActor(request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%d", id)
}

func (h *Handler) putActor(w http.ResponseWriter, r *http.Request, id int) {
	var request actor.ActorRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	actor, err := h.services.PutActor(id, request)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) patchActor(w http.ResponseWriter, r *http.Request, id int) {
	var request actor.ActorRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	actor, err := h.services.PatchActor(id, request)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(actor)
	fmt.Fprintf(w, "%s", reqBodyBytes.String())
}

func (h *Handler) deleteActor(w http.ResponseWriter, id int) {
	err := h.services.DeleteActor(id)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
}
