package controllers

import (
	"os"
	"errors"
	"strings"
	"net/http"
	"api/utils"
	"api/models"
	util "api/utils/mail"
	"github.com/rs/zerolog"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Define Logger Instance
var logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()

// Request for reset token
func ForgotPassword(c *gin.Context) {

	// get user's email from JSON request
	var userCredential *models.ForgotPasswordInput

	if err := c.ShouldBindJSON(&userCredential); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	// Check if user with the email exists
	var user models.User
	if err := models.DB.First(&user, "email_address = ?", strings.ToLower(userCredential.EmailAddress)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
            c.IndentedJSON(http.StatusNotFound, gin.H{
                "status":  "fail",
                "message": "User not found",
            })
        } else {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{
                "status":  "fail",
                "message": "Failed to retrieve user",
            })
        }
        return
	}

	// Check if user is verified
	if !user.IsEmailVerified {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"message": "Account not verified, please verify your email",
		})
		return
	}

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error().Err(errors.New("reading environment variables failed")).Msgf("%v",err)
	}

	// Generate password reset Code
	resetToken, err := utils.GenerateRandomString(20)

	// Encode generated reset code
	passwordResetToken := util.Encode(resetToken)

	// Update user password reset data in Database 
	passwordReset := models.PasswordReset{}
	passwordReset.Username = user.Username
	passwordReset.EmailAddress = user.EmailAddress
	passwordReset.PasswordResetCode = passwordResetToken
	models.DB.Save(&passwordReset) // -> Save data to DB

	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	var firstName = user.FirstName

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Email Data
	emailData := util.EmailData{
		URL:       os.Getenv("CLIENT_ORIGIN") + "/api/auth/reset-password/" + resetToken,
		FirstName: firstName,
		Subject:   "Your Loyalty Points Program password reset token",
	}

	// Send Password Reset Email
	util.SendEmail(&user, &emailData, "passwordReset.html")

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"status": "fail",
			"message": "There was an error sending email",
		})
		return
	}

	message := "We sent an email with a reset password code to " + user.EmailAddress

	// Success
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": message,
	})
}

// Reset user's password
func ResetPassword(c *gin.Context) {

	// Get resetToken from parameter
	resetToken := c.Params.ByName("resetToken")

	// Get user password reset data
	var userCredential *models.ResetPasswordInput

	if err := c.ShouldBindJSON(&userCredential); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": err.Error(),
		})
		return
	}

	// Check if passwords match
	if userCredential.NewPassword != userCredential.ConfirmPassword {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Passwords do not match",
		})
		return
	}

	// Hash the new password
	hashedPassword, _ := models.HashPassword(userCredential.NewPassword)

	passwordResetToken := util.Encode(resetToken)

	// Update the user's password in the database
	var updatedUser models.PasswordReset
	var user models.User

	result := models.DB.First(&updatedUser, "password_reset_code = ?", passwordResetToken)

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

	user.Password = hashedPassword
	updatedUser.PasswordResetCode = ""

	models.DB.Save(&user)
	models.DB.Save(&updatedUser)

	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Password reset successful"},
	)
}
