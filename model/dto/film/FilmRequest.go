package film

import (
	"encoding/json"
	"time"
)

const dateFormat = "2006-01-02"

type FilmRequest struct {
	Name        *string    `json:"name" validate:"min=1,max=150"`
	Description *string    `json:"description" validate:"max=1000"`
	ReleaseDate *time.Time `json:"releaseDate"`
	Rating      *int       `json:"rating" validate:"min=0,max=10"`
	ActorsId    *[]int     `json:"actorsId"`
}

func (film *FilmRequest) UnmarshalJSON(p []byte) error {
	var aux struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		ReleaseDate *string `json:"releaseDate"`
		Rating      *int    `json:"rating"`
		ActorsId    *[]int  `json:"actorsId"`
	}

	err := json.Unmarshal(p, &aux)
	if err != nil {
		return err
	}
	var t time.Time

	if aux.ReleaseDate != nil {
		t, err = time.Parse(dateFormat, *aux.ReleaseDate)
		if err != nil {
			return err
		}
	}

	film.Name = aux.Name
	film.ReleaseDate = &t
	film.Description = aux.Description
	film.Rating = aux.Rating
	film.ActorsId = aux.ActorsId

	return nil
}
