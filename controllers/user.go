/*
This file contains logic for various API functions
*/

package controllers

import (
	"net/http"
	"api/utils/token"
	models "api/models"
	"github.com/gin-gonic/gin"
)


// create a new user
func CreateAccount(c *gin.Context){
	
	// Validate Input
	var input models.CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// // Create User
	user := models.User{}
	
	user.Username = input.Username
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Gender = input.Gender
	user.EmailAddress = input.EmailAddress
	user.Password = input.Password
	user.City = input.City
	user.PhoneNumber = input.PhoneNumber

	_,err := user.SaveUser()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
    }
	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "User Created!",
	})
}

// login user
func Login(c *gin.Context){

	// Validate input
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Login User
	user := models.User{}
	user.Username = input.Username
	user.Password = input.Password

	token, err := models.LoginCheck(user.Username, user.Password)

	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Username or Password is incorrect",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Login Success",
		"token": token,
	})
}

// get current user
func GetCurrentUser(c *gin.Context){

	uid, err := token.ExtractTokenID(c)

	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user,err := models.GetUser(uid)

	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Success!",
		"data": user,
	})
}
