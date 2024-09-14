package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"main/models"
	"regexp"
)

var AuthTokenRegex = regexp.MustCompile("Token (\\S+)")

var (
	NoAuthTokenError        = errors.New("no auth token")
	InvalidTokenFormatError = errors.New("invalid token format")
)

func CreateAuthToken() *models.AuthToken {
	rawToken := make([]byte, 32)
	_, err := rand.Read(rawToken)
	if err != nil {
		panic(err)
	}
	token := models.AuthToken(base64.URLEncoding.EncodeToString(rawToken))
	return &token
}

func GetTokenFromRequest(context *gin.Context) (*models.AuthToken, error) {
	authorization := context.GetHeader("Authorization")
	if len(authorization) == 0 {
		return nil, NoAuthTokenError
	}

	match := AuthTokenRegex.FindStringSubmatch(authorization)
	if len(match) == 0 {
		return nil, InvalidTokenFormatError
	}
	token := models.AuthToken(match[1])
	return &token, nil
}
