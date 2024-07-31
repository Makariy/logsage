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

func RenderAuthorizationHeader(token AuthToken) string {
	return fmt.Sprintf("Token %s", token)
}

func Authorize(context *gin.Context, user *models.User) (AuthToken, error) {
	token := CreateAuthToken()
	context.Header("Authorization", RenderAuthorizationHeader(token))

	return token, SetUserByToken(user, token)
}

func SignUpUser(context *gin.Context, email, pass string) (*models.User, AuthToken, error) {
	user, err := repository.CreateUser(email, pass)
	if err != nil {
		return nil, "", err
	}
	token, err := Authorize(context, user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func LogoutUser(context *gin.Context) error {
	context.Header("Authorization", "")
	ctxUser, _ := context.Get("user")
	user := ctxUser.(*models.User)
	return DelUser(user)
}
