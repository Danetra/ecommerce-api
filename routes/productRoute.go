package routes

import (
	"ecommerce-api/controllers"

	"github.com/gin-gonic/gin"
)

func productRoute(api *gin.RouterGroup) {
	api.GET("/products", controllers.GetProducts)
	api.POST("/products", controllers.CreateProduct)
	api.GET("/products/:id", controllers.GetProductByID)
	api.PUT("/products/:id", controllers.UpdateProduct)
	api.DELETE("/products/:id", controllers.DeleteProduct)
}
