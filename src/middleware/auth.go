package middleware

import (
	"github.com/gin-gonic/gin"
	"main/auth"
	"main/models"
	"net/http"
)

func LoginRequired(ctx *gin.Context) {
	user, err := GetUserFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "Forbidden",
		})
		ctx.Abort()
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}

func GetUserFromRequest(ctx *gin.Context) (*models.User, error) {
	ctxUser, exists := ctx.Get(UserKey)
	if exists {
		user := ctxUser.(*models.User)
		return user, nil
	}

	token, err := auth.GetTokenFromRequest(ctx)
	if err != nil {
		return nil, err
	}

	user, err := auth.GetUserByToken(models.AuthToken(token))
	if err != nil {
		return nil, err
	}
	return user, nil
}
