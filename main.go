package main

import (
	"ecommerce-api/config"
	"ecommerce-api/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ecommerce-api/docs"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	config.DbConfig()

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(
		config.DB,
		"postgres",
		migrations,
		migrate.Up,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Applied %d migrations\n", n)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Static("/uploads", "./uploads")

	routes.ApiRoutes(router)
	fmt.Println("Server running at http://localhost:8080")
	router.Run(":8080")

}
