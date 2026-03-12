package models

import "time"

type User struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username" validate:"required,min=2,max=255"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8,max=255"`
	Role      string    `json:"role" validate:"oneof=admin user"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`

	Token        string `json:"token"`
	Refreshtoken string `json:"refreshtoken"`

	FavouriteGeneres []Genre `json:"favouritegeneres"`
}
