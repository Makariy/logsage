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

func AddStatsRoutes(router *gin.Engine) {
	group := router.Group("/stats")
	group.Use(middleware.LoginRequired)

	group.GET("/all/", middleware.AttachDateRange, handleGetTotalStats)
	group.GET("/category/:id/",
		middleware.AttachDateRange,
		middleware.AttachUserAndModel[models.Category](),
		handleGetCategoryStats)
	group.GET("/account/:id/",
		middleware.AttachDateRange,
		middleware.AttachUserAndModel[models.Account](),
		handleGetAccountStats)
}

func handleGetCategoryStats(ctx *gin.Context) {
	category, exists := middleware.GetFromContext[*models.Category](ctx, middleware.ModelKey)
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	dateRange, exists := middleware.GetFromContext[*forms.DateTimeRange](ctx, middleware.DateRangeKey)
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
		DateRange:       dateRange.ToDateRange(),
		Category: forms.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Type: category.Type,
		},
		Stats: *stats,
	})
}

func handleGetAccountStats(ctx *gin.Context) {
	account, exists := middleware.GetFromContext[*models.Account](ctx, middleware.ModelKey)
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	dateRange, exists := middleware.GetFromContext[*forms.DateTimeRange](ctx, middleware.DateRangeKey)
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
		DateRange:       dateRange.ToDateRange(),
		Account:         *accountForm,
		Stats:           *stats,
	})
}

func handleGetTotalStats(ctx *gin.Context) {
	user, err := middleware.GetUserFromRequest(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	dateRange, exists := middleware.GetFromContext[*forms.DateTimeRange](ctx, middleware.DateRangeKey)
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
		DateRange:       dateRange.ToDateRange(),
		Stats:           *stats,
	})
}
