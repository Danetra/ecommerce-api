package main

import (
	"ecommerce-api/config"
	"ecommerce-api/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
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

	routes.ApiRoutes(router)
	fmt.Println("Server running at http://localhost:8080")
	router.Run(":8080")

}
