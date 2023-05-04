package controllers

import (
	"net/http"
	models "api/models"
	"github.com/gin-gonic/gin"
)


// POST user
// create a new user
func AddUser(c *gin.Context){

	db, err := models.ConnectDB()

	if err != nil {
		panic(err)
	}
	
	// Validate Input
	// var input models.CreateUserInput
	// if err := c.BindJSON(&input); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{
	// 		"error" : err.Error(),
	// 	})
	// 	return 
	// }

	// Create User
	sqlStatement := `INSERT INTO Users (
		username, firstname, lastname, gender, email, password, phone_number
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`
    _, err = db.Exec(sqlStatement, "Njoguu", "Alan", "Njogu","male", "smith@acme.com", "alannjoguu", "+09088789")
    if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		panic(err)
    } else {
        c.IndentedJSON(http.StatusCreated, gin.H{
			"message": "User Created!",
		})
    }
	return
}
