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
		middleware.AttachRelativeCurrencyByDefaultKeys(),
		middleware.AttachDateRangeByDefaultKeys(),
		handleGetTotalCategoriesStats)
	group.GET(
		"/account/all/",
		middleware.AttachDateRangeByDefaultKeys(),
		handleGetTotalAccountsStats)
	group.GET("/category/:id/",
		middleware.AttachRelativeCurrencyByDefaultKeys(),
		middleware.AttachDateRangeByDefaultKeys(),
		middleware.AttachUserAndModelByDefaultKeys[models.Category](),
		handleGetCategoryStats)
	group.GET("/account/:id/",
		middleware.AttachDateRangeByDefaultKeys(),
		middleware.AttachUserAndModelByDefaultKeys[models.Account](),
		handleGetAccountStats)
	group.GET("/interval/",
		middleware.AttachRelativeCurrencyByDefaultKeys(),
		middleware.AttachDateRangeByDefaultKeys(),
		middleware.AttachTimeIntervalByDefaultKeys(),
		handleGetTimeIntervalStats)
}

func handleGetCategoryStats(ctx *gin.Context) {
	category := middleware.GetFromContext[models.Category](ctx, middleware.ModelKey)
	relativeCurrency := middleware.GetFromContext[models.Currency](
		ctx,
		middleware.RelativeCurrencyKey,
	)
	dateRange := middleware.GetDateRangeFromContext(ctx)

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
	account := middleware.GetFromContext[models.Account](ctx, middleware.ModelKey)
	dateRange := middleware.GetDateRangeFromContext(ctx)

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
	user, _ := middleware.GetUserFromRequest(ctx)
	dateRange := middleware.GetDateRangeFromContext(ctx)

	relativeCurrency := middleware.GetFromContext[models.Currency](
		ctx,
		middleware.RelativeCurrencyKey,
	)

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
	user, _ := middleware.GetUserFromRequest(ctx)
	dateRange := middleware.GetDateRangeFromContext(ctx)

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

func handleGetTimeIntervalStats(ctx *gin.Context) {
	user, _ := middleware.GetUserFromRequest(ctx)
	dateRange := middleware.GetDateRangeFromContext(ctx)
	relativeCurrency := middleware.GetFromContext[models.Currency](
		ctx, middleware.RelativeCurrencyKey,
	)
	timeInterval := middleware.GetFromContext[forms.TimeInterval](ctx, middleware.TimeIntervalKey)

	stats, err := repository.GetTimeIntervalStats(
		user.ID,
		dateRange.FromDate,
		dateRange.ToDate,
		models.TimeIntervalStep(timeInterval.Interval),
		relativeCurrency,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(
			500,
			forms.ErrorResponse{
				Status: "error",
				Error:  err.Error(),
			},
		)
		return
	}
	ctx.JSON(http.StatusOK, forms.TimeIntervalStatsResponse{
		SuccessResponse: forms.Success,
		DateRange:       dateRange.ToDateRange(),
		Stats:           *stats,
	})
}
