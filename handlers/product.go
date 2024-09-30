package handlers

import (
	"net/http"
	"strconv"

	"github.com/NikosDouras/simpler-project-nik/database"
	"github.com/NikosDouras/simpler-project-nik/models"
	"github.com/gin-gonic/gin"
)

// CreateProduct
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (name, price, quantity) are required and must be valid"})
		return
	}

	result := database.DB.Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProducts with pagination
func GetProducts(c *gin.Context) {
	var products []models.Product
	var total int64

	limitParam := c.DefaultQuery("limit", "10")
	pageParam := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	result := database.DB.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	database.DB.Model(&models.Product{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetProduct by ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")

	// Validate that the ID is numeric
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates a product by ID
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Validate that the ID is numeric
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var existingProduct models.Product

	result := database.DB.First(&existingProduct, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Create a struct for binding update data
	var updateData struct {
		Name        string  `json:"name" binding:"required"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Quantity    int     `json:"quantity" binding:"required,gte=0"`
		Description string  `json:"description"` // Optional field
	}

	// Bind incoming JSON data to updateData struct
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (name, price, quantity) are required and must be valid"})
		return
	}

	// Update existing product with new data
	existingProduct.Name = updateData.Name
	existingProduct.Price = updateData.Price
	existingProduct.Quantity = updateData.Quantity
	existingProduct.Description = updateData.Description

	// Save changes to database
	database.DB.Save(&existingProduct)

	c.JSON(http.StatusOK, existingProduct)
}

// DeleteProduct deletes a product by ID
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Validate that the ID is numeric
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	database.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
