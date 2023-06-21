/*
This file contains logic for various API functions concerning 
Membership Service and Auth
*/

package controllers

import (
	"os"
	"time"
	"errors"
	"strings"
	"strconv"
	"net/http"
	"api/utils"
	models "api/models"
	mail "api/utils/mail"
	token "api/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thanhpk/randstr"
)

// create a new user
func CreateAccount(c *gin.Context){
	
	// Validate Input
	var input models.CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": err.Error(),
		})
		return
	}

	// Hash Password
	hashedPassword, err := models.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "error",
			"message": err.Error()})
		return
	}

	// // Create User
	user := models.User{}
	
	user.Username = input.Username
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Gender = input.Gender
	user.EmailAddress = input.EmailAddress
	user.Password = hashedPassword
	user.City = input.City
	user.PhoneNumber = input.PhoneNumber
	user.RedeemablePoints = 150	// allocate 150 points by default 
	user.VirtualCardNumber = utils.GenerateCardNumber()

	_,errr := user.SaveUser()

	if errr != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"status": "fail", 
			"message": "User with that credential already exists",
		})
		return
	} 

	// Generate Verification Code
	code := randstr.String(20)

	verification_code := mail.Encode(code)

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

	// Save to Transactions Table
	models.SaveToTransactions(user.Username, "EARN", "WELCOME BONUS", 150, "")

	// ? Send Email
	emailData := mail.EmailData{
		URL:       os.Getenv("CLIENT_ORIGIN") + "/api/auth/verify-email/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	mail.SendEmail(&user, &emailData, "verificationCode.html")

	message := "We sent an email with a verification code to " + user.EmailAddress

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success", 
		"message": message,
	})
}

// verify Email
func  VerifyEmail(c *gin.Context) {

	code := c.Params.ByName("secret_code")
	verification_code := mail.Encode(code)

	var updatedUser models.VerifyEmails
	var user models.User

	result := models.DB.First(&updatedUser, "secret_code = ?", verification_code)
	
	res := models.DB.First(&user, "username = ?", updatedUser.Username)

	if res.Error != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": 
			"Invalid verification code or user doesn't exist",
		})
		return
	}

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail", 
			"message": "Invalid verification code or user doesn't exist",
		})
		return
	}

	if user.IsEmailVerified {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"status": "fail",
			"message": "User already verified",
		})
		return
	}

	user.IsEmailVerified = true
	updatedUser.SecretCode = ""
	user.RedeemablePoints = user.RedeemablePoints + 350		// Allocate 350 points on email verfification

	models.DB.Save(&user)
	models.DB.Save(&updatedUser)

	// Save to Transactions Table
	models.SaveToTransactions(user.Username, "EARN", "CONFIRM EMAIL ADDRESS BONUS", 350, "")

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Email verified successfully",
	})
}

// login user
func Login(c *gin.Context) {

	// Validate Input
	var input *models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": err.Error(),
		})
		return
	}
	
	// Get User data of user logging in
	var user models.User
	result := models.DB.First(&user, "email_address = ?", strings.ToLower(input.EmailAddress))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Invalid email or Password",
		})
		return
	}

	// Check if User is verified
	if !user.IsEmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Please verify your email"})
		return
	}

	// Verify password given
	if err := models.VerifyPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error().Err(errors.New("reading environment variables failed")).Msgf("%v",err)
	}

	// Generate Token
	tokenExpiresInStr := os.Getenv("TOKEN_EXPIRES_IN")
	TOKEN_MAXAGEStr := os.Getenv("TOKEN_MAXAGE")

	TOKEN_MAXAGE, err := strconv.Atoi(TOKEN_MAXAGEStr)
	if err != nil{
		return
	}

	TOKEN_EXPIRES_IN, err := time.ParseDuration(tokenExpiresInStr)
	if err != nil {
		logger.Error().Err(errors.New("invalid value for TOKEN_EXPIRES_IN")).Msgf("%v",err)
	}
	
	// Generate Token
	token, err := token.GenerateToken(TOKEN_EXPIRES_IN, user.ID, os.Getenv("TOKEN_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("token", token, TOKEN_MAXAGE*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token": token,
	})
}

// logout user
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "You have been logged out!",
	})
}

// Get User Data
func GetCurrentUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		Username: currentUser.Username,
		FirstName: currentUser.FirstName,
		LastName: currentUser.LastName,
		EmailAddress: currentUser.EmailAddress, 
		PhoneNumber: currentUser.PhoneNumber,
		RedeemablePoints: currentUser.RedeemablePoints,
		City: currentUser.City,	
		VirtualCardNumber: currentUser.VirtualCardNumber,
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": userResponse,
		},
	})
}

// Change Password
func ChangePassword(c *gin.Context){
	// Get Logged In user
	currentUser := c.MustGet("currentUser").(models.User)

	// Get the new password and confirm password from the request body
	var passwordReset struct {
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&passwordReset); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": err.Error(),
		})
		return
	}

	// Check if the new password and confirm password match
	if passwordReset.NewPassword != passwordReset.ConfirmPassword {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "Passwords do not match",
		})
		return
	}

	// Hash the new password
	hashedPassword, err := models.HashPassword(passwordReset.NewPassword)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "Failed to generate hashed password",
		})
		return
	}

	// Update the user's password in the database
	err = models.UpdateUserPassword(currentUser.ID, hashedPassword)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "Failed to update password",
		})
		return
	}

	// Password reset successful
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Password change successful",
	})
}
