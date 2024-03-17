package presenter

type ActorResponse struct {
	Id       int    `json:"id""`
	Name     string `json:"name""`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
	FilmsId  []int  `json:"filmsId"`
}
