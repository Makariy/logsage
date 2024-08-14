package main

import (
	"github.com/gin-gonic/gin"
	_ "main/cache"
	_ "main/db_connector"
	"main/routes"
)

func main() {
	router := gin.Default()

	routes.SetupAllRoutes(router)

	err := router.Run("localhost:8000")
	if err != nil {
		panic("Got an error running gin: " + err.Error())
	}
}
