/*
This file contains logic for various API functions
*/

package controllers

import (
	"os"
	"strings"
	"net/http"
	"api/utils/mail"
	"api/utils/token"
	models "api/models"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
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

	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		c.IndentedJSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that credential already exists"})
		return
	} 

	// Generate Verification Code
	code := randstr.String(20)

	verification_code := utils.Encode(code)

	// Update User verification data in Database
	verifyEmails := models.VerifyEmails{}
	verifyEmails.Username = user.Username
	verifyEmails.EmailAddress = user.EmailAddress
	verifyEmails.SecretCode = verification_code
	models.DB.Save(&verifyEmails)

	var firstName = user.FirstName

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       os.Getenv("CLIENT_ORIGIN") + "/api/auth/verify-email/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(&user, &emailData)

	message := "We sent an email with a verification code to " + user.EmailAddress

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success", 
		"message": message,
	})
}

// verify Email
func  VerifyEmail(c *gin.Context) {

	code := c.Params.ByName("secret_code")
	verification_code := utils.Encode(code)

	var updatedUser models.VerifyEmails
	var user models.User

	result := models.DB.First(&updatedUser, "secret_code = ?", verification_code)
	
	res := models.DB.First(&user, "username = ?", updatedUser.Username)

	if res.Error != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if user.IsEmailVerified {
		c.IndentedJSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	user.IsEmailVerified = true
	updatedUser.SecretCode = ""
	
	models.DB.Save(&user)
	models.DB.Save(&updatedUser)

	// Allocate free starter points on account creation
	models.AllocatePoints(updatedUser.ID, updatedUser.Username)

	c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
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

// User to view their points balance
func GetPointsData(c *gin.Context){

	uid, err := token.ExtractTokenID(c)

	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userData,err := models.GetPointsByID(uid)

	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Success!",
		"data": userData,
	})
}
