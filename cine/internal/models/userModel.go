package models

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username" validate:"required,min=2,max=255"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=255"`
	Role      string `json:"role" validate:"oneof=admin user"`
	CreatedAt string `json:"createdat"`
	UpdatedAt string `json:"updatedat"`

	Token string
}
