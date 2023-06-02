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
    Count int   `json:"count"`
}
