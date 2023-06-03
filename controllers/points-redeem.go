package controllers

import (
    "net/http"
	models "api/models"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
)

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
    
    // Calculate the total price for redemption
    totalPoints := selectedProduct.Price * request.Quantity

    // TODO: Deduct request.Quantity from selectedProduct.Count

    // Check if the user has enough points for redemption
    if  user.RedeemablePoints < totalPoints {
        c.IndentedJSON(http.StatusBadRequest, gin.H{
            "status": "fail",
            "message": "Insufficient points",
        })
        return
    }

    // Deduct points from the user's account
    user.RedeemablePoints = user.RedeemablePoints - totalPoints

    if err := models.DB.Save(&user).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{
            "status":  "fail",
            "message": "Failed to deduct points",
        })
        return
    }

    // TODO: Create a new transaction record
    models.SaveToTransactions(user.Username, "REDEEM", "Product redemption",totalPoints)

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
