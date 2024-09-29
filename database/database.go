package database

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/NikosDouras/simpler-project-nik/models"
)

var DB *gorm.DB

// Connect initializes the database connection depending on the environment.
func Connect() {
	var err error
	var dialector gorm.Dialector

	// Check if running tests
	if os.Getenv("GO_ENV") == "test" {
		// Use in-memory SQLite for testing
		fmt.Println("Using in-memory SQLite for testing")
		dialector = sqlite.Open(":memory:")
	} else {
		// Use PostgreSQL for non-test environments
		fmt.Println("Using PostgreSQL database")
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
		dialector = postgres.Open(dsn)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database! ", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database! ", err)
	}
}
