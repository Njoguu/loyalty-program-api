/*
This file contains logic for various API functions concerning 
member transaction history data.
*/

package controllers

import (
	"time"
	"errors"
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
                "message": "transactions not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "error",
                "message": "failed to retrieve transactions",
            })
        }
        return
	}

	// Convert date and time format
	for i := range userTransactions {
		parsedDate, err := time.Parse(time.RFC3339, userTransactions[i].Date)
		if err != nil {
			logger.Error().Err(errors.New("error parsing date")).Msgf("%v",err)
			return
		}
		userTransactions[i].Date = parsedDate.Format("January 2, 2006")
		
		// time field is in "2006-01-02T15:04:05Z" format
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z", userTransactions[i].Time)
		if err != nil {
			logger.Error().Err(errors.New("error parsing time")).Msgf("%v",err)
			return
		}
		userTransactions[i].Time = parsedTime.Format("15:04:05")
	}

	
	
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Transaction History Retrieval successful!",
		"data": gin.H{
			"transactions": userTransactions,
		},
	})
}
