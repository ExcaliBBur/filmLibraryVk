package pkg

import (
	"log"
	"net/http"
)

func JWTAuthAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s request on %s", r.Method, r.RequestURI)

		err := ValidateJWT(w, r)
		if err != nil {
			log.Printf("Invalid JWT token")
			http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
			return
		}
		error := ValidateAdminRoleJWT(w, r)
		if error != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func JWTAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s request on %s", r.Method, r.RequestURI)

		err := ValidateJWT(w, r)
		if err != nil {
			log.Printf("Invalid JWT token")
			http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
			return
		}
		error := ValidateUserRoleJWT(w, r)
		if error != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
