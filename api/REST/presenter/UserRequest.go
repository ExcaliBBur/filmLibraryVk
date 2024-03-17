package presenter

type UserRequest struct {
	Username *string `json:"username" validate:"min=2"`
	Password *string `json:"password" validate:"min=8,max=20"`
	Role     *string `json:"role"`
}
