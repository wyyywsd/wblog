package middleware

import "github.com/gin-gonic/gin"

// SetMiddleware 设置中间件
func SetMiddleware(key string, value interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(key, value)
		context.Next()
	}
}
