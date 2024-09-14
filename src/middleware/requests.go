package middleware

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"main/repository"
	"main/utils"
	"net/http"
)

func AttachUser(ctxKey string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user, err := GetUserFromRequest(ctx)
		if err != nil {
			ctx.Abort()
			return
		}
		ctx.Set(ctxKey, user)
	}
}
func AttachUserByDefaultKeys() func(ctx *gin.Context) {
	return AttachUser(UserKey)
}

func AttachModelID(urlKey, ctxKey string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		modelID, err := utils.ShouldParseID(urlKey, ctx)
		if err != nil {
			ctx.Abort()
			return
		}

		ctx.Set(ctxKey, &modelID)
	}
}

func AttachModelIDByDefaultKeys() func(*gin.Context) {
	return AttachModelID(ModelIdURLKey, ModelIdKey)
}

func AttachModel[Model models.UserGettable](
	modelIdKey string,
	ctxKey string,
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		modelID, exists := ShouldGetFromContext[*models.ModelID](ctx, modelIdKey)
		if !exists {
			ctx.Abort()
			return
		}

		model, err := repository.GetModelByID[Model](*modelID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, forms.ErrorResponse{
				Status: "Not found",
				Error:  "No item found with this ID",
			})
			return
		}

		ctx.Set(ctxKey, model)
	}
}

func AttachModelByDefaultKeys[Model models.UserGettable]() func(ctx *gin.Context) {
	return AttachModel[Model](ModelIdKey, ModelKey)
}

func AttachDateRange(ctxKey string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		dateForm, err := utils.ShouldGetQuery[forms.DateRange](ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "Invalid date",
			})
			return
		}
		if dateForm.ToDate-dateForm.FromDate <= 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "Invalid date range",
			})
			return
		}

		ctx.Set(ctxKey, dateForm.ToDateTimeRange())
	}
}

func AttachDateRangeByDefaultKeys() func(ctx *gin.Context) {
	return AttachDateRange(DateRangeKey)
}

func AttachRelativeCurrency(ctxKey string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		relativeCurrency, err := utils.ShouldGetQuery[forms.RelativeCurrency](ctx)
		if err != nil || relativeCurrency.Symbol == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "Invalid relative currency",
			})
			return
		}

		currency, err := repository.GetCurrencyBySymbol(relativeCurrency.Symbol)
		if err != nil || currency == nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "Could not find relative currency",
			})
			return
		}

		ctx.Set(ctxKey, currency)
	}
}

func AttachRelativeCurrencyByDefaultKeys() func(ctx *gin.Context) {
	return AttachRelativeCurrency(RelativeCurrencyKey)
}

func AttachTimeInterval(ctxKey string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		timeInterval, err := utils.ShouldGetQuery[forms.TimeInterval](ctx)
		if err != nil || timeInterval.Interval == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "Invalid time interval",
			})
			return
		}

		if !models.IsTimeIntervalDefined(timeInterval.Interval) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "Time interval is not defined",
			})
			return
		}

		ctx.Set(ctxKey, timeInterval)
	}
}

func AttachTimeIntervalByDefaultKeys() func(ctx *gin.Context) {
	return AttachTimeInterval(TimeIntervalKey)
}
