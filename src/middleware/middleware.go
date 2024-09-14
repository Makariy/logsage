package middleware

import "github.com/gin-gonic/gin"

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
