package routes

import (
	"ecommerce-api/controllers"

	"github.com/gin-gonic/gin"
)

func transactionRoute(api *gin.RouterGroup) {
	api.GET("/transactions", controllers.GetTransactions)
	api.POST("/transactions", controllers.CreateTransaction)
	api.GET("/transactions/history", controllers.TransactionHistory)
	api.POST("/transactions/:id/payment", controllers.TransactionPayment)
}
