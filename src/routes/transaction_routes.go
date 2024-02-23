package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"main/utils"
)

func AddTransactionRoutes(router *gin.Engine) {
	group := router.Group("/transaction")
	group.Use(utils.LoginRequired)

	group.GET("/all/", handleGetUserModels[models.Transaction, forms.TransactionResponse, forms.TransactionsResponse])
	group.GET("/get/:id/", modelPermissionRequired[models.Transaction], handleGetUserModel[models.Transaction, forms.TransactionResponse])
	group.POST("/create/", handleCreateModel[models.Transaction, forms.TransactionForm, forms.TransactionResponse])
	group.PATCH("/patch/:id/", modelPermissionRequired[models.Transaction], handlePatchModel[models.Transaction, forms.TransactionForm, forms.TransactionResponse])
	group.DELETE("/delete/:id/", modelPermissionRequired[models.Transaction], handleDeleteModel[models.Transaction, forms.TransactionResponse])
}
