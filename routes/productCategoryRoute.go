package routes

import (
	"ecommerce-api/controllers"

	"github.com/gin-gonic/gin"
)

func productCategoryRoute(api *gin.RouterGroup) {
	api.GET("/product-categories", controllers.GetProductCategories)
	api.POST("/product-categories", controllers.CreateProductCategory)
	api.GET("/product-categories/:id", controllers.GetProductCategoryByID)
	api.PUT("/product-categories/:id", controllers.UpdateProductCategory)
	api.DELETE("/product-categories/:id", controllers.DeleteProductCategory)
}
