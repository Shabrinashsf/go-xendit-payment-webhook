package middleware

import (
	"net/http"
	"os"

	"github.com/Shabrinashsf/go-xendit-payment-webhook/utils/response"
	"github.com/gin-gonic/gin"
)

func Xendit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("X-CALLBACK-TOKEN")
		if authHeader == "" {
			res := response.BuildResponseFailed("Failed to get callback token", "Missing X-CALLBACK-TOKEN header", nil)
			ctx.JSON(http.StatusUnauthorized, res)
			ctx.Abort()
			return
		}
		xenditKey := os.Getenv("XENDIT_WEBHOOK_TOKEN")
		if authHeader != xenditKey {
			res := response.BuildResponseFailed("Unauthorized token", "Invalid X-CALLBACK-TOKEN", nil)
			ctx.JSON(http.StatusUnauthorized, res)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
