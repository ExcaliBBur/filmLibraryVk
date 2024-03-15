package actor

import (
	"encoding/json"
	"time"
)

const dateFormat = "2006-01-02"

type ActorRequest struct {
	Name     *string    `json:"name""`
	Sex      *string    `json:"sex"`
	Birthday *time.Time `json:"birthday"`
	FilmsId  *[]int     `json:"filmsId"`
}

func (actor *ActorRequest) UnmarshalJSON(p []byte) error {
	var aux struct {
		Name     *string `json:"name"`
		Sex      *string `json:"sex"`
		Birthday *string `json:"birthday"`
		FilmsId  *[]int  `json:"filmsId"`
	}

	err := json.Unmarshal(p, &aux)
	if err != nil {
		return err
	}
	var t time.Time

	if aux.Birthday != nil {
		t, err = time.Parse(dateFormat, *aux.Birthday)
		if err != nil {
			return err
		}
	}

	actor.Name = aux.Name
	actor.Sex = aux.Sex
	actor.Birthday = &t
	actor.FilmsId = aux.FilmsId

	return nil
}
