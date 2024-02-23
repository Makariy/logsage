package models

type UserProtected interface {
	GetUser() *User
}
