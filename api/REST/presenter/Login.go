package presenter

type Login struct {
	Username string `json:"username" validate:"min=2"`
	Password string `json:"password" validate:"min=8,max=20"`
}
