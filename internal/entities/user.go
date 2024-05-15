package entities

type User struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
