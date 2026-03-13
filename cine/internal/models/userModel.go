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

	Favouritegeneres []Genre `json:"favouritegeneres"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type FoundUser struct {
	UserId           string  `json:"userid"`
	Username         string  `json:"username"`
	Email            string  `json:"email"`
	Password         string  `json:"password"`
	Role             string  `json:"role"`
	Token            string  `json:"token"`
	Refreshtoken     string  `json:"refreshtoken"`
	Favouritegeneres []Genre `json:"favouritegeneres"`
}
