package utils

import (
	"github.com/gin-gonic/gin"
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
