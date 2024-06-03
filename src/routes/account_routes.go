package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/middleware"
	"main/models"
	"main/repository"
	"main/utils"
)

func AddAccountRoutes(router *gin.Engine) {
	group := router.Group("/account")
	group.Use(middleware.LoginRequired)

	group.GET("/all/",
		handleGetUserModels[models.Account, forms.AccountResponse, forms.AccountsResponse](
			repository.GetUserModels[models.Account],
		),
	)
	group.GET("/get/:id/",
		middleware.AttachUserAndModel[models.Account](),
		handleGetUserModel[models.Account, forms.AccountResponse](repository.GetModelByID[models.Account]))
	group.POST("/create/",
		handleCreateUserModel[models.Account, forms.AccountForm, forms.AccountResponse](
			utils.ShouldGetForm[forms.AccountForm], repository.CreateModel[models.Account]))
	group.PATCH("/patch/:id/",
		middleware.AttachUserAndModel[models.Account](),
		handlePatchModel[models.Account, forms.AccountForm, forms.AccountResponse](
			utils.ShouldGetForm[forms.AccountForm], repository.PatchModel[models.Account]))
	group.DELETE("/delete/:id/",
		middleware.AttachUserAndModel[models.Account](),
		handleDeleteModel[models.Account, forms.AccountResponse](repository.DeleteModel[models.Account]))
}
