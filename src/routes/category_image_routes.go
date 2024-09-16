package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"main/repository"
)

func AddCategoryImageRoutes(route *gin.Engine) {
	group := route.Group("/category_image")

	group.GET(
		"/all/",
		handleGetAllModels[
			models.CategoryImage,
			forms.CategoryImageResponse,
			forms.CategoryImagesResponse,
		](repository.GetAllModels[models.CategoryImage]),
	)
}
