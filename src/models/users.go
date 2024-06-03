package models

import (
	"time"
)

type User struct {
	ID           ModelID   `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Email        string    `gorm:"column:email;unique"`
	PasswordHash string    `gorm:"column:password_hash"`
	LastLogin    time.Time `gorm:"column:last_login;default:current_timestamp"`
}

func (_ *User) TableName() string {
	return "auth_user"
}
