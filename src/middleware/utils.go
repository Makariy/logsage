package middleware

import "github.com/gin-gonic/gin"

func GetFromContext[Item any](ctx *gin.Context, key string) (Item, bool) {
	item, exists := ctx.Get(key)
	if !exists {
		var emptyItem Item
		return emptyItem, false
	}

	var result Item
	result = item.(Item)
	return result, true
}

func GroupMiddlewares(middlewares []gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, middleware := range middlewares {
			middleware(ctx)
			if ctx.IsAborted() {
				return
			}
		}
	}
}
