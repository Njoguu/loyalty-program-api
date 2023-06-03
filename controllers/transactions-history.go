/*
This file contains logic for various API functions concerning 
member transaction history data.
*/

package controllers

import (
    "net/http"
	models "api/models"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
)


// Function to view a user's transaction history
func ViewTransactions(c *gin.Context){
	// Get the logged-in user from the authentication middleware
	user := c.MustGet("currentUser").(models.User)

	// Retrieve Transactions from Database
	var userTransactions []models.Transactions
	if err := models.DB.Where("username = ?", user.Username).Find(&userTransactions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
            c.IndentedJSON(http.StatusNotFound, gin.H{
                "status":  "fail",
                "message": "Transactions not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "fail",
                "message": "Failed to retrieve transactions",
            })
        }
        return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Transaction History Retrieval successful!",
		"data": gin.H{
			"transactions": userTransactions,	// TODO: Return all Transactions as List
		},
	})
}
