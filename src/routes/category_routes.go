package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"main/utils"
)

func AddCategoryRoutes(router *gin.Engine) {
	group := router.Group("/category")
	group.Use(utils.LoginRequired)

	group.POST("/create/", handleCreateModel[models.Category, forms.CategoryForm, forms.CategoryResponse])
	group.GET("/all/", handleGetUserModels[models.Category, forms.CategoryResponse, forms.CategoriesResponse])
	group.GET("/get/:id/", modelPermissionRequired[models.Category], handleGetUserModel[models.Category, forms.CategoryResponse])
	group.PATCH("/patch/:id/", modelPermissionRequired[models.Category], handlePatchModel[models.Category, forms.CategoryForm, forms.CategoryResponse])
	group.DELETE("/delete/:id/", modelPermissionRequired[models.Category], handleDeleteModel[models.Category, forms.CategoryResponse])
}
