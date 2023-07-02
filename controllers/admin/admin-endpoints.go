package controllers

import (
    "net/http"
	models "api/models"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
)


// Add product
func AddProduct(c *gin.Context) {
	// Get the product ID and quantity from the request body
    var newProduct models.Products

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": err.Error(),
		})
		return
	}

	// Create a new product with custom values
	product := models.Products{
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Quantity:       newProduct.Quantity,
	}

	models.SaveProduct(&product)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"status": "success",
		"message": "Product added successfully",
	})
}

// Update product
func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	var updatedProduct models.Products
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": err.Error(),
		})
		return
	}

	var product models.Products
	if err := models.DB.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
            c.IndentedJSON(http.StatusNotFound, gin.H{
                "status":  "fail",
                "message": "Product not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "fail",
                "message": "Failed to retrieve product",
            })
        }
        return
	}

	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	product.Quantity = updatedProduct.Quantity

	if err := models.DB.Save(&product).Error; err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Failed to update product",
		})
		return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Product updated successfully",
	})
}

// Delete product
func DeleteProduct(c *gin.Context){
	// Get product id from parameter
	productId := c.Param("id")

	var product models.Products
	if err := models.DB.First(&product, productId).Error; err != nil{
		if err == gorm.ErrRecordNotFound {
            c.IndentedJSON(http.StatusNotFound, gin.H{
                "status":  "fail",
                "message": "Product not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "err",
                "message": "Failed to delete product",
            })
        }
        return	
	}

	if err := models.DB.Delete(&product).Error; err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Failed to delete product",
		})
		return
	}
	
	c.IndentedJSON(http.StatusNoContent, gin.H{
		"status": "success",
		"message": "Product deleted successfully",
	})
}

// Get List of all Products - Admin
func GetAllProducts(c *gin.Context){

	// Retrieve Products from Database
	var allProducts []models.Products
	if err := models.DB.Find(&allProducts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
            c.IndentedJSON(http.StatusNotFound, gin.H{
                "status":  "fail",
                "message": "Products not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "fail",
                "message": "Failed to retrieve products",
            })
        }
        return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Products Retrieval successful!",
		"data": gin.H{
			"products": allProducts,	
		},
	})
}
