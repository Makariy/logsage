package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"main/db_connector"
	"main/models"
)

func HashPassword(password string) string {
	sha := sha256.New()
	sha.Write([]byte(password))
	hashSum := sha.Sum(nil)
	return hex.EncodeToString(hashSum)
}

func GetUserByEmail(email string) (*models.User, error) {
	db := db_connector.GetConnection()

	var user models.User
	tx := db.Find(&user, "email = ?", email)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func CreateUser(email, password string) (*models.User, error) {
	passwordHash := HashPassword(password)
	user := models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	return CreateModel(&user)
}

func DeleteUser(id models.ModelID) (*models.User, error) {
	return DeleteModel[models.User](id)
}
