package utils

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"net/http"
	"strconv"
)

func ShouldGetForm[T any](ctx *gin.Context) (*T, error) {
	form, err := forms.ValidateForm[T](ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  err.Error(),
		})
		return nil, err
	}
	return form, nil
}

func ShouldGetQuery[T any](ctx *gin.Context) (*T, error) {
	form, err := forms.ValidateQuery[T](ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  err.Error(),
		})
		return nil, err
	}
	return form, nil
}

func ShouldParseID(ctx *gin.Context) (models.ModelID, error) {
	strID := ctx.Param("id")
	if strID == "" {
		strID = ctx.Query("id")
	}
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  "Could not parse id",
		})
		return 0, err
	}
	return models.ModelID(id), nil
}
