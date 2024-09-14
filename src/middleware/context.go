package middleware

import (
	"github.com/gin-gonic/gin"
	"main/forms"
)

func ShouldGetFromContext[Item any](ctx *gin.Context, key string) (Item, bool) {
	item, exists := ctx.Get(key)
	if !exists {
		var emptyItem Item
		return emptyItem, false
	}

	var result Item
	result = item.(Item)
	return result, true
}

func GetFromContext[T any](ctx *gin.Context, key string) *T {
	model, exists := ShouldGetFromContext[*T](ctx, key)
	if !exists {
		ctx.AbortWithStatus(500)
		return nil
	}
	return model
}

func GetDateRangeFromContext(ctx *gin.Context) *forms.DateTimeRange {
	dateRange, exists := ShouldGetFromContext[*forms.DateTimeRange](ctx, DateRangeKey)
	if !exists {
		ctx.AbortWithStatus(500)
		return nil
	}
	return dateRange
}
