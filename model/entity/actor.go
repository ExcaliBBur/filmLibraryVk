package entity

import "time"

type Actor struct {
	Id       int       `json:"id""`
	Name     string    `json:"name""`
	Sex      string    `json:"sex"`
	Birthday time.Time `json:"birthday" time:"2006-01-02"`
}
