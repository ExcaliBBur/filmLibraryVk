package entity

import "time"

type Film struct {
	Id          int       `json:"id"`
	Name        string    `json:"name" validate:"min=1,max=150"`
	Description string    `json:"description" validate:"max=1000"`
	ReleaseDate time.Time `json:"releaseDate"`
	Rating      int       `json:"rating" validate:"min=0,max=10"`
}
