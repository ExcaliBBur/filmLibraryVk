package entity

import "time"

type Film struct {
	Id          int       `json:"id"`
	Name        string    `json:"name" validate:"minLength:1, maxLength: 150"`
	Description string    `json:"description validate:"maxLength: 1000""`
	ReleaseDate time.Time `json"releaseDate"`
	Rating      int       `json:"rating validate:"min:0, max:10""`
}
