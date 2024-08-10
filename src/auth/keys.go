package auth

import (
	"fmt"
	"main/models"
)

func getUserKeyPatternById(id models.ModelID) string {
	return fmt.Sprintf("user_%d_*", id)
}

func getUserKeyPatternByToken(token models.AuthToken) string {
	return fmt.Sprintf("user_*_%s", string(token))
}

func createUserKeyByToken(user *models.User, token models.AuthToken) string {
	return fmt.Sprintf("user_%d_%s", user.ID, string(token))
}
