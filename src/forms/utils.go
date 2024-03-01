package forms

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateForm[T any](ctx *gin.Context) (*T, error) {
	var form T
	err := ctx.BindJSON(&form)
	if err != nil {
		return nil, err
	}

	err = v.Struct(&form)
	if err != nil {
		return nil, err
	}
	return &form, nil
}
