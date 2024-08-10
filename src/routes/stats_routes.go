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

	group.GET(
		"/category/all/",
		middleware.AttachRelativeCurrency,
		middleware.AttachDateRange,
		handleGetTotalCategoriesStats,
	)
	group.GET("/account/all/", middleware.AttachDateRange, handleGetTotalAccountsStats)
	group.GET("/category/:id/",
		middleware.AttachRelativeCurrency,
		middleware.AttachDateRange,
		middleware.AttachUserAndModel[models.Category](),
		handleGetCategoryStats,
	)
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

	relativeCurrency, exists := middleware.GetFromContext[*models.Currency](
		ctx, middleware.RelativeCurrencyKey,
	)
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	stats, err := repository.GetCategoryStats(
		category.ID,
		dateRange.FromDate,
		dateRange.ToDate,
		relativeCurrency,
	)
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

func handleGetTotalCategoriesStats(ctx *gin.Context) {
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

	relativeCurrency, exists := middleware.GetFromContext[*models.Currency](
		ctx, middleware.RelativeCurrencyKey,
	)
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	stats, err := repository.GetTotalCategoriesStats(
		user.ID,
		dateRange.FromDate,
		dateRange.ToDate,
		relativeCurrency,
	)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, forms.TotalCategoriesStatsResponse{
		SuccessResponse: forms.Success,
		DateRange:       dateRange.ToDateRange(),
		Stats:           *stats,
	})
}

func handleGetTotalAccountsStats(ctx *gin.Context) {
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

	stats, err := repository.GetTotalAccountsStats(user.ID, dateRange.FromDate, dateRange.ToDate)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, forms.TotalAccountsStatsResponse{
		SuccessResponse: forms.Success,
		DateRange:       dateRange.ToDateRange(),
		Stats:           *stats,
	})
}
