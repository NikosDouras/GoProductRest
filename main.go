package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/NikosDouras/simpler-project-nik/database"
	"github.com/NikosDouras/simpler-project-nik/handlers"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database.Connect()

	// Initialize Gin router
	r := gin.Default()

	// Routes
	r.POST("/products", handlers.CreateProduct)
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProduct)
	r.PUT("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)

	// Run the server
	r.Run(":8080") // Default port is 8080
}
