package forms

import (
	"main/models"
	"time"
)

type UserForm struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=6,max=64"`
}

type UserResponse struct {
	*SuccessResponse
	Email            string           `json:"email"`
	Token            models.AuthToken `json:"token"`
	LastLogin        time.Time        `json:"lastLogin"`
	RegistrationDate time.Time        `json:"registrationDate"`
}
