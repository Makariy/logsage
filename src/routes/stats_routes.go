package routes

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"main/repository"
	"main/utils"
	"net/http"
)

func AddStatsRoutes(router *gin.Engine) {
	group := router.Group("/stats")
	group.Use(utils.LoginRequired)

	group.GET("/all/", includeDateRange, handleGetTotalStats)
	group.GET("/category/:id/", includeDateRange, modelPermissionRequired[models.Category], handleGetCategoryStats)
	group.GET("/account/:id/", includeDateRange, modelPermissionRequired[models.Account], handleGetAccountStats)
}

const DateRangeKey = "dateRange"

func includeDateRange(ctx *gin.Context) {
	dateForm, err := utils.ShouldGetForm[forms.DateRange](ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  "Invalid date",
		})
		return
	}
	if dateForm.FromDate.Compare(dateForm.ToDate) > 0 {
		ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  "Invalid date range",
		})
		return
	}

	ctx.Set(DateRangeKey, dateForm)
}

func handleGetCategoryStats(ctx *gin.Context) {
	category, exists := utils.GetFromContext[models.Category](ctx, "item")
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	dateRange, exists := utils.GetFromContext[forms.DateRange](ctx, DateRangeKey)
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	stats, err := repository.GetCategoryStats(category.ID, dateRange.FromDate, dateRange.ToDate)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, forms.CategoryStatsResponse{
		SuccessResponse: forms.Success,
		Category: forms.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Type: category.Type,
		},
		Stats: *stats,
	})
}

func handleGetAccountStats(ctx *gin.Context) {
	account, exists := utils.GetFromContext[models.Account](ctx, "item")
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	dateRange, exists := utils.GetFromContext[forms.DateRange](ctx, DateRangeKey)
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	stats, err := repository.GetAccountStats(account.ID, dateRange.FromDate, dateRange.ToDate)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	accountForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](account)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, forms.AccountStatsResponse{
		SuccessResponse: forms.Success,
		Account:         *accountForm,
		Stats:           *stats,
	})
}

func handleGetTotalStats(ctx *gin.Context) {
	user, err := utils.GetUserFromRequest(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	dateRange, exists := utils.GetFromContext[forms.DateRange](ctx, DateRangeKey)
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	stats, err := repository.GetTotalStats(user.ID, dateRange.FromDate, dateRange.ToDate)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, forms.TotalStatsResponse{
		SuccessResponse: forms.Success,
		Stats:           *stats,
	})
}
