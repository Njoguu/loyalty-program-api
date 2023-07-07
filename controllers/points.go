package controllers

import (
    "net/http"
	models "api/models"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
)

// Function to redeem points
func RedeemPoints(c *gin.Context){
	// Get the logged-in user from the authentication middleware
    user := c.MustGet("currentUser").(models.User)

    // Get the product ID and quantity from the request body
    var request struct {
        ProductName  string  `json:"product_name"`
        Quantity int `json:"product_quantity"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": err.Error(), 
		})
        return
    }
    
	// Retrieve the product from the database
    var selectedProduct models.Products
    if err := models.DB.Where("name = ?", request.ProductName).First(&selectedProduct).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.IndentedJSON(http.StatusNotFound, gin.H{
                "status":  "fail",
                "message": "product not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "error",
                "message": "failed to retrieve product",
            })
        }
        return
    }
    
    // Calculate the total price for redemption
    totalPoints := selectedProduct.Price * request.Quantity

    if selectedProduct.Quantity < 0 {
        // Handle the error condition where the user tries to redeem more quantity than available
        c.IndentedJSON(http.StatusBadRequest, gin.H{
            "status": "fail",
            "message": "insufficient quantity available for redemption",
        })
        return
    }

    // Deduct request.Quantity from selectedProduct.Count
    selectedProduct.Quantity -= request.Quantity
    models.DB.Save(&selectedProduct).Where("quantity > ?", 0)

    // Check if the user has enough points for redemption
    if  user.RedeemablePoints < totalPoints {
        c.IndentedJSON(http.StatusBadRequest, gin.H{
            "status": "fail",
            "message": "insufficient points",
        })
        return
    }

    // Deduct points from the user's account
    user.RedeemablePoints = user.RedeemablePoints - totalPoints

    if err := models.DB.Save(&user).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to update points",
        })
        return
    }

    // Create a new transaction record
    if err := models.SaveToTransactions(user.Username, "REDEEM", "PRODUCT REDEMPTION",totalPoints, selectedProduct.Name); err != nil{
        c.IndentedJSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "message": err.Error(),
        })
    }

    // Return the redeemed product and remaining points to the user
    c.IndentedJSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "Redemption successful",
        "redeemed_product": selectedProduct.Name,
        "redeemed_points": totalPoints,
        "balance": gin.H{
            "redeemable_points": user.RedeemablePoints,
        },
    })
}

// TODO: Use Virtual Card Number to allocate earned points from purchase

// TODO: Add the transaction to user's transaction history

// Calculate the number of points to award
func CalculatePoints(amount float64) int {
	
	if amount < 500 {
		points := int(amount / 100) * 1
		return points
	} else if amount <= 1500 {
		points := int(amount / 100) * 4
		return points
	} else if amount <= 3000 {
		return 70
	} else if amount <= 5000 {
		return 150
	} else {
		return 300
	}

}