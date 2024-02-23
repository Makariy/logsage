package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"main/utils"
)

func AddAccountRoutes(router *gin.Engine) {
	group := router.Group("/account")
	group.Use(utils.LoginRequired)

	group.GET("/all/", handleGetUserModels[models.Account, forms.AccountResponse, forms.AccountsResponse])
	group.GET("/get/:id/", modelPermissionRequired[models.Account], handleGetUserModel[models.Account, forms.AccountResponse])
	group.POST("/create/", handleCreateModel[models.Account, forms.AccountForm, forms.AccountResponse])
	group.PATCH("/patch/:id/", modelPermissionRequired[models.Account], handlePatchModel[models.Account, forms.AccountForm, forms.AccountResponse])
	group.DELETE("/delete/:id/", modelPermissionRequired[models.Account], handleDeleteModel[models.Account, forms.AccountResponse])
}
