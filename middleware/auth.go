package middleware

import (
	"github.com/gin-gonic/gin"
)

func UUIDRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.GetHeader("Authorization")
		ctx.Set("uuid", uuid)
	}
}
