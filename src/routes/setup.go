package routes

import "github.com/gin-gonic/gin"

func SetupAllRoutes(router *gin.Engine) {
	AddAuthRoutes(router)
	AddAccountRoutes(router)
	AddCurrencyRoutes(router)
	AddCategoryRoutes(router)
	AddCategoryImageRoutes(router)
	AddStatsRoutes(router)
	AddTransactionRoutes(router)
}
