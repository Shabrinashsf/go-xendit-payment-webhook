package routes

import (
	"github.com/Shabrinashsf/go-xendit-payment-webhook/controller"
	"github.com/Shabrinashsf/go-xendit-payment-webhook/middleware"
	"github.com/gin-gonic/gin"
)

func Transaction(route *gin.Engine, transactionController controller.TransactionController) {
	routes := route.Group("/transaction")
	{
		routes.POST("/buy", transactionController.CreateTransaction)
		routes.POST("/webhook/xendit", middleware.Xendit(), transactionController.XenditWebhook)
	}
}
