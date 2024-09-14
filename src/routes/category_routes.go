package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/middleware"
	"main/models"
	"main/repository"
	"main/utils"
)

func AddCategoryRoutes(router *gin.Engine) {
	group := router.Group("/category")
	group.Use(middleware.LoginRequired)

	group.POST(
		"/create/",
		handleCreateUserModel[
			models.Category,
			forms.CategoryForm,
			forms.CategoryResponse,
		](utils.ShouldGetForm[forms.CategoryForm], repository.CreateModel[models.Category]),
	)
	group.GET("/all/",
		handleGetUserModels[models.Category, forms.CategoryResponse, forms.CategoriesResponse](
			repository.GetUserModels[models.Category],
		),
	)
	group.GET("/get/:id/",
		middleware.AttachUserAndModelByDefaultKeys[models.Category](),
		handleGetUserModel[models.Category, forms.CategoryResponse](repository.GetModelByID[models.Category]),
	)
	group.PATCH(
		"/patch/:id/",
		middleware.AttachUserAndModelByDefaultKeys[models.Category](),
		handlePatchModel[
			models.Category,
			forms.CategoryForm,
			forms.CategoryResponse,
		](utils.ShouldGetForm[forms.CategoryForm], repository.PatchModel[models.Category]))
	group.DELETE(
		"/delete/:id/",
		middleware.AttachUserAndModelByDefaultKeys[models.Category](),
		handleDeleteModel[models.Category, forms.CategoryResponse](repository.DeleteModel[models.Category]))
}
