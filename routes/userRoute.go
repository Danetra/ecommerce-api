package routes

import (
	"ecommerce-api/controllers"

	"github.com/gin-gonic/gin"
)

func userRoute(api *gin.RouterGroup) {
	api.GET("/users/:id", controllers.GetUserByID)
	//api.PUT("/users/:id", controllers.UpdateUser)
	//api.DELETE("/users/:id", controllers.DeleteRole)
}
