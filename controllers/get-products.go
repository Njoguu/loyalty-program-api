package controllers

import (
    "net/http"
	models "api/models"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
)

// Product list users can see
type PublicProduct struct {
	Name string `json:"product_name"`
	Description string `json:"product_description"`
	Price int `json:"product_price"`
}

// Get List of all Products - User
func GetProducts(c *gin.Context) {
	var products []models.Products

	if err := models.DB.Find(&products).Error; err != nil{
		if err == gorm.ErrRecordNotFound {
            c.IndentedJSON(http.StatusNotFound, gin.H{
                "status":  "fail",
                "message": "products not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "fail",
                "message": "failed to retrieve products",
            })
        }
        return
	}

	publicProducts := make([]PublicProduct, len(products))
	for i, product := range products {
		publicProducts[i] = PublicProduct{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "products retrieval successful",
		"data": gin.H{
			"products": publicProducts,
		},
	})
}

