package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"main/models"
	"main/repository"
)

func GetUserByCredentials(email, password string) (*models.User, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.PasswordHash != repository.HashPassword(password) {
		return nil, errors.New("passwords do not match")
	}

	return user, nil
}

func RenderAuthorizationHeader(token *models.AuthToken) string {
	return fmt.Sprintf("Token %s", *token)
}

func getOrCreateTokenForUser(user *models.User) (*models.AuthToken, error) {
	_, err := GetUserByID(user.ID)
	if err == nil {
		return GetTokenByUserID(user.ID)
	}

	token := CreateAuthToken()
	err = SetUserByToken(user, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func Authorize(context *gin.Context, user *models.User) (*models.AuthToken, error) {
	token, err := getOrCreateTokenForUser(user)
	if err != nil {
		_ = DelUser(user)
		return nil, err
	}
	context.Header("Authorization", RenderAuthorizationHeader(token))

	err = repository.UpdateLastLogin(user)
	if err != nil {
		_ = DelUser(user)
		return nil, err
	}

	return token, err
}

func SignUpUser(context *gin.Context, email, pass string) (*models.User, *models.AuthToken, error) {
	user, err := repository.CreateUser(email, pass)
	if err != nil {
		return nil, nil, err
	}
	token, err := Authorize(context, user)
	if err != nil {
		return nil, nil, err
	}
	return user, token, nil
}

func LogoutUser(context *gin.Context) error {
	context.Header("Authorization", "")
	ctxUser, _ := context.Get("user")
	user := ctxUser.(*models.User)
	return DelUser(user)
}
