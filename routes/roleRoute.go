package routes

import (
	"ecommerce-api/controllers"

	"github.com/gin-gonic/gin"
)

func roleRoute(api *gin.RouterGroup) {
	api.GET("/roles", controllers.GetRoles)
	api.POST("/roles", controllers.CreateRole)
	api.GET("/roles/:id", controllers.GetRoleByID)
	api.PUT("/roles/:id", controllers.UpdateRole)
	api.DELETE("/roles/:id", controllers.DeleteRole)
}
