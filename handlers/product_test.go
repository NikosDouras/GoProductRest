// handlers/product_test.go

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/NikosDouras/simpler-project-nik/database"
	"github.com/NikosDouras/simpler-project-nik/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

// TestMain sets up the environment and runs all tests
func TestMain(m *testing.M) {
	// Set the environment to "test" so the database.Connect() uses SQLite
	os.Setenv("GO_ENV", "test")

	// Connect to the test database
	database.Connect()

	// Set up the router once for all tests
	router = setupRouter()

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}

// setupTestDB clears the database before each test
func setupTestDB() {
	// Clear all records before each test
	database.DB.Exec("DELETE FROM products")
}

// setupRouter initializes the router with routes for testing
func setupRouter() *gin.Engine {
	r := gin.Default()

	// Register routes for testing
	r.POST("/products", CreateProduct)
	r.GET("/products", GetProducts)
	r.GET("/products/:id", GetProduct)
	r.PUT("/products/:id", UpdateProduct)
	r.DELETE("/products/:id", DeleteProduct)

	return r
}

func TestCreateProduct(t *testing.T) {
	// Clear the database for isolated tests
	setupTestDB()

	product := models.Product{
		Name:        "Test Product",
		Price:       10.99,
		Description: "A test product",
		Quantity:    100,
	}

	productJSON, _ := json.Marshal(product)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdProduct models.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	assert.Nil(t, err)
	assert.Equal(t, product.Name, createdProduct.Name)
	assert.Equal(t, product.Price, createdProduct.Price)
}

func TestCreateProductValidationErrors(t *testing.T) {
	setupTestDB()
	// Test missing fields
	invalidProduct := map[string]interface{}{
		"Price": 10.99, // Missing "Name" and "Quantity"
	}
	productJSON, _ := json.Marshal(invalidProduct)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response, "error") // Expect an error message
}

func TestGetProducts(t *testing.T) {
	// Clear the database for isolated tests
	setupTestDB()

	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response, "data")
}

func TestGetNonExistentProduct(t *testing.T) {
	// Clear the database for isolated tests
	setupTestDB()

	req, _ := http.NewRequest("GET", "/products/9999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Product not found", response["error"])
}

func TestGetProductInvalidID(t *testing.T) {
	setupTestDB()

	// Request with a non-existent product ID
	req, _ := http.NewRequest("GET", "/products/9999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Product not found", response["error"])
}

func TestGetProductsPagination(t *testing.T) {
	setupTestDB()

	// Create multiple products
	for i := 1; i <= 15; i++ {
		product := models.Product{
			Name:        fmt.Sprintf("Product %d", i),
			Price:       float64(i * 10),
			Description: fmt.Sprintf("Description for product %d", i),
			Quantity:    i * 5,
		}
		database.DB.Create(&product)
	}

	// Request with pagination
	req, _ := http.NewRequest("GET", "/products?limit=5&page=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Verify pagination details
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "page")
	assert.Contains(t, response, "limit")
	assert.Equal(t, float64(5), response["limit"]) // Ensure limit is correct
	assert.Equal(t, float64(2), response["page"])  // Ensure page is correct
}

func TestGetProductsDefaultPagination(t *testing.T) {
	setupTestDB()

	// Create a few products
	for i := 1; i <= 20; i++ {
		product := models.Product{
			Name:        fmt.Sprintf("Product %d", i),
			Price:       float64(i * 10),
			Description: fmt.Sprintf("Description for product %d", i),
			Quantity:    i * 5,
		}
		database.DB.Create(&product)
	}

	// Request without limit and page parameters
	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Verify defaults
	assert.Equal(t, float64(10), response["limit"]) // Default limit = 10
	assert.Equal(t, float64(1), response["page"])   // Default page = 1
}

func TestGetProductsPaginationEdgeCases(t *testing.T) {
	setupTestDB()

	for i := 1; i <= 30; i++ {
		product := models.Product{
			Name:        fmt.Sprintf("Product %d", i),
			Price:       float64(i * 10),
			Description: fmt.Sprintf("Description for product %d", i),
			Quantity:    i * 5,
		}
		database.DB.Create(&product)
	}

	// Test with limit=0 (should fallback to default limit)
	req, _ := http.NewRequest("GET", "/products?limit=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(10), response["limit"])

	// Test with page=1000 (should return no products)
	req, _ = http.NewRequest("GET", "/products?page=1000", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Empty(t, response["data"])
}

func TestUpdateProduct(t *testing.T) {
	// Clear the database for isolated tests
	setupTestDB()

	// First, create a product to update
	product := models.Product{
		Name:        "Original Product",
		Price:       15.00,
		Description: "A product to be updated",
		Quantity:    20,
	}

	database.DB.Create(&product)

	// Update request
	updatedProduct := models.Product{
		Name:        "Updated Product",
		Price:       20.00,
		Description: "An updated product",
		Quantity:    50,
	}

	productJSON, _ := json.Marshal(updatedProduct)
	req, _ := http.NewRequest("PUT", "/products/"+strconv.Itoa(int(product.ID)), bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.Product
	err := json.Unmarshal(w.Body.Bytes(), &updated)
	assert.Nil(t, err)
	assert.Equal(t, updatedProduct.Name, updated.Name)
	assert.Equal(t, updatedProduct.Price, updated.Price)
}

func TestUpdateNonExistentProduct(t *testing.T) {
	setupTestDB()

	updatedProduct := models.Product{
		Name:        "Updated Product",
		Price:       20.00,
		Description: "An updated product",
		Quantity:    50,
	}

	productJSON, _ := json.Marshal(updatedProduct)
	req, _ := http.NewRequest("PUT", "/products/9999", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Product not found", response["error"])
}

func TestUpdateProductIncompleteFields(t *testing.T) {
	setupTestDB()

	// Create a product to update
	product := models.Product{
		Name:        "Original Product",
		Price:       15.00,
		Description: "A product to be updated",
		Quantity:    20,
	}

	database.DB.Create(&product)

	// Missing "Name" and "Quantity"
	updatedProduct := map[string]interface{}{
		"Price": 20.00,
	}

	productJSON, _ := json.Marshal(updatedProduct)
	req, _ := http.NewRequest("PUT", "/products/"+strconv.Itoa(int(product.ID)), bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response, "error") // Expect an error message about missing fields
}

func TestUpdateProductInvalidID(t *testing.T) {
	setupTestDB()

	updatedProduct := map[string]interface{}{
		"Name":     "Updated Product",
		"Price":    20.00,
		"Quantity": 10,
	}

	productJSON, _ := json.Marshal(updatedProduct)

	// Request with a non-existent product ID
	req, _ := http.NewRequest("PUT", "/products/9999", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Product not found", response["error"])
}

func TestDeleteProduct(t *testing.T) {
	// Clear the database for isolated tests
	setupTestDB()

	// First, create a product to delete
	product := models.Product{
		Name:        "Product to Delete",
		Price:       12.50,
		Description: "This product will be deleted",
		Quantity:    30,
	}

	database.DB.Create(&product)

	// Delete request
	req, _ := http.NewRequest("DELETE", "/products/"+strconv.Itoa(int(product.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Product deleted successfully", response["message"])

	// Try to get the deleted product
	req, _ = http.NewRequest("GET", "/products/"+strconv.Itoa(int(product.ID)), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteProductInvalidID(t *testing.T) {
	setupTestDB()

	req, _ := http.NewRequest("DELETE", "/products/invalid-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response, "error")
}

func TestDeleteNonExistentProduct(t *testing.T) {
	// Clear the database for isolated tests
	setupTestDB()

	req, _ := http.NewRequest("DELETE", "/products/9999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Product not found", response["error"])
}
