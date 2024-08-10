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

func CreateAuthToken() models.AuthToken {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return models.AuthToken(base64.URLEncoding.EncodeToString(token))
}

func GetTokenFromRequest(context *gin.Context) ([]byte, error) {
	authorization := context.GetHeader("Authorization")
	if len(authorization) == 0 {
		return nil, NoAuthTokenError
	}

	match := AuthTokenRegex.FindStringSubmatch(authorization)
	if len(match) == 0 {
		return nil, InvalidTokenFormatError
	}
	return []byte(match[1]), nil
}
