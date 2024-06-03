package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/middleware"
	"main/models"
	"main/repository"
	"main/utils"
	"net/http"
)

func AddTransactionRoutes(router *gin.Engine) {
	group := router.Group("/transaction")
	group.Use(middleware.LoginRequired)

	group.GET(
		"/all/",
		middleware.AttachDateRange,
		handleGetTransactions,
	)
	group.GET("/get/:id/",
		middleware.AttachUserAndModel[models.Transaction](),
		handleGetUserModel[models.Transaction, forms.TransactionResponse](
			repository.GetModelByID[models.Transaction],
		),
	)
	group.POST("/create/",
		handleCreateUserModel[models.Transaction, forms.TransactionForm, forms.TransactionResponse](
			utils.ShouldGetForm[forms.TransactionForm], repository.CreateModel[models.Transaction],
		),
	)
	group.PATCH("/patch/:id/",
		middleware.AttachUserAndModel[models.Transaction](),
		handlePatchModel[models.Transaction, forms.TransactionForm, forms.TransactionResponse](
			utils.ShouldGetForm[forms.TransactionForm], repository.PatchModel[models.Transaction],
		),
	)
	group.DELETE("/delete/:id/",
		middleware.AttachUserAndModel[models.Transaction](),
		handleDeleteModel[models.Transaction, forms.TransactionResponse](
			repository.DeleteModel[models.Transaction],
		),
	)
}

func handleGetTransactions(ctx *gin.Context) {
	dateRange, exists := middleware.GetFromContext[*forms.DateTimeRange](ctx, middleware.DateRangeKey)
	if !exists {
		return
	}

	user, err := middleware.GetUserFromRequest(ctx)
	if err != nil {
		return
	}

	transactions, err := repository.GetUserModelsWithDateRange[models.Transaction](
		user.ID,
		dateRange.FromDate,
		dateRange.ToDate,
	)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, forms.ErrorResponse{
			Status: "Internal Server Error",
			Error:  err.Error(),
		})
	}

	response, err := renderResponses[models.Transaction, forms.TransactionResponse, forms.TransactionsResponse](
		transactions,
	)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
