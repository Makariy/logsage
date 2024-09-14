package auth

import (
	"errors"
	"fmt"
	"main/models"
	"regexp"
)

func getUserKeyPatternById(id models.ModelID) string {
	return fmt.Sprintf("user_%d_*", id)
}

func getUserKeyPatternByToken(token *models.AuthToken) string {
	return fmt.Sprintf("user_*_%s", string(*token))
}

func createUserKeyByToken(user *models.User, token *models.AuthToken) string {
	return fmt.Sprintf("user_%d_%s", user.ID, string(*token))
}

func parseTokenByKey(key string) (*models.AuthToken, error) {
	re, _ := regexp.Compile("user_\\d+_([\\s\\S]+)")
	match := re.FindStringSubmatch(key)
	if len(match) == 0 {
		return nil, errors.New("cannot match token")
	}
	token := models.AuthToken(match[1])
	return &token, nil
}
