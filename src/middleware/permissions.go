package middleware

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"net/http"
)

func CheckUserModelPermission[Model models.UserGettable](userKey, modelKey string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user, exists := ShouldGetFromContext[*models.User](ctx, userKey)
		if !exists {
			ctx.Abort()
			return
		}

		model, exists := ShouldGetFromContext[*Model](ctx, modelKey)
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
}

func CheckUserModelPermissionByDefaultKeys[Model models.UserGettable]() func(ctx *gin.Context) {
	return CheckUserModelPermission[Model](UserKey, ModelKey)
}

func AttachUserAndModel[Model models.UserGettable](
	userKey string,
	modelIDURLKey string,
	modelIDCtxKey string,
	modelKey string,
) gin.HandlerFunc {
	return GroupMiddlewares([]gin.HandlerFunc{
		AttachUser(userKey),
		AttachModelID(modelIDURLKey, modelIDCtxKey),
		AttachModel[Model](modelIDCtxKey, modelKey),
		CheckUserModelPermission[Model](userKey, modelKey),
	})
}

func AttachUserAndModelByDefaultKeys[Model models.UserGettable]() gin.HandlerFunc {
	return AttachUserAndModel[Model](
		UserKey,
		ModelIdURLKey,
		ModelIdKey,
		ModelKey,
	)
}
