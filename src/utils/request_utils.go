package utils

import (
	"github.com/gin-gonic/gin"
	"main/auth"
	"main/forms"
	"main/models"
	"net/http"
	"strconv"
)

func GetUserFromRequest(ctx *gin.Context) (*models.User, error) {
	ctxUser, exists := ctx.Get("user")
	if exists {
		user := ctxUser.(*models.User)
		return user, nil
	}

	token, err := auth.GetTokenFromRequest(ctx)
	if err != nil {
		return nil, err
	}

	user, err := auth.GetUserByToken(auth.AuthToken(token))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetFromContext[T any](ctx *gin.Context, key string) (*T, bool) {
	value, exists := ctx.Get(key)
	if !exists {
		return nil, exists
	}
	result, ok := value.(T)
	if !ok {
		return nil, false
	}
	return &result, true
}

func ShouldGetForm[T any](ctx *gin.Context) (*T, error) {
	form, err := forms.ValidateForm[T](ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  err.Error(),
		})
		return nil, err
	}
	return form, nil
}

func ShouldParseID(ctx *gin.Context) (uint, error) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  "Could not parse id",
		})
		return 0, err
	}
	return uint(id), nil
}
