package film

type FilmResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ReleaseDate string `json"releaseDate"`
	Rating      int    `json:"rating`
	ActorsId    []int  `json:"actorsId"`
}
