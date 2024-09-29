// models/product.go

package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name" binding:"required"`       // Add validation tag
	Price       float64   `json:"price" binding:"required,gt=0"` // Price must be greater than 0
	Description string    `json:"description"`
	Quantity    int       `json:"quantity" binding:"required,gte=0"` // Quantity must be non-negative
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
