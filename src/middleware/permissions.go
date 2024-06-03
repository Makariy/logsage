package middleware

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"net/http"
)

func UserModelPermission[Model models.UserProtected](ctx *gin.Context) {
	user, exists := GetFromContext[*models.User](ctx, UserKey)
	if !exists {
		ctx.Abort()
		return
	}

	model, exists := GetFromContext[*Model](ctx, ModelKey)
	if !exists {
		ctx.Abort()
		return
	}

	if (*model).GetUser().ID != user.ID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, forms.ErrorResponse{
			Status: "Forbidden",
			Error:  "Operations on this model are forbidden",
		})
		return
	}
}

func AttachUserAndModel[Model models.UserProtected]() gin.HandlerFunc {
	return GroupMiddlewares([]gin.HandlerFunc{
		AttachUser,
		AttachModelID,
		AttachModel[Model],
		UserModelPermission[Model],
	})
}
