package models

import (
	"github.com/jinzhu/gorm"
)

// Product table - details about a product
type Products struct {
    gorm.Model
    Name string `json:"product_name"`
    Description string  `json:"product_description"`
    Price int   `json:"product_price"`
    Quantity int   `json:"quantity"`
}

// MISC FUNCTIONS
// Function to save product added by admin
func SaveProduct(product *Products){
    // Save Product to DB
    DB.Create(&product)
}