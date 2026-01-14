package routes

import (
	"ecommerce-api/controllers"
	"ecommerce-api/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	api := router.Group("/api")

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	// Add middleware JWT for Auth
	api.Use(middleware.JWTAuth())

	roleRoute(api)
	userRoute(api)
}
