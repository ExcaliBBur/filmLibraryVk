package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handleError(w http.ResponseWriter, err error, status int) {
	log.Printf("Error: %s", err.Error())
	http.Error(w, err.Error(), status)
}



func getPathId(w http.ResponseWriter, r *http.Request, prefix string) int {
	idString := strings.TrimPrefix(r.URL.Path, prefix)
	id, err := strconv.Atoi(idString)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return -1
	}
	return id
}
