package middleware

import (
	"github.com/gin-gonic/gin"
	"main/forms"
	"main/models"
	"main/repository"
	"main/utils"
	"net/http"
)

const (
	UserKey      = "user"
	ModelKey     = "model"
	ModelIdKey   = "modelId"
	DateRangeKey = "dateRange"
)

func AttachUser(ctx *gin.Context) {
	user, err := GetUserFromRequest(ctx)
	if err != nil {
		ctx.Abort()
		return
	}
	ctx.Set(UserKey, user)
}

func AttachModelID(ctx *gin.Context) {
	modelID, err := utils.ShouldParseID(ctx)
	if err != nil {
		ctx.Abort()
		return
	}

	ctx.Set(ModelIdKey, &modelID)
}

func AttachModel[Model models.UserGettable](ctx *gin.Context) {
	modelID, exists := GetFromContext[*models.ModelID](ctx, ModelIdKey)
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

	ctx.Set(ModelKey, model)
}

func AttachDateRange(ctx *gin.Context) {
	dateForm, err := utils.ShouldGetQuery[forms.DateRange](ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  "Invalid date",
		})
		return
	}
	if dateForm.ToDate-dateForm.FromDate < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  "Invalid date range",
		})
		return
	}

	ctx.Set(DateRangeKey, dateForm.ToDateTimeRange())
}
