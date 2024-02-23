package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
)

func AddCurrencyRoutes(router *gin.Engine) {
	group := router.Group("/currency")

	group.GET("/all/", handleGetAllModels[models.Currency, forms.CurrencyResponse, forms.CurrenciesResponse])
}
