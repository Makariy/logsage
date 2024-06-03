package main

import (
	"github.com/gin-gonic/gin"
	_ "main/cache"
	_ "main/db_connector"
	"main/routes"
)

func main() {
	router := gin.Default()

	routes.AddAuthRoutes(router)
	routes.AddAccountRoutes(router)
	routes.AddCurrencyRoutes(router)
	routes.AddCategoryRoutes(router)
	routes.AddStatsRoutes(router)
	routes.AddTransactionRoutes(router)

	err := router.Run("localhost:8000")
	if err != nil {
		panic("Got an error executing gin")
	}
}
