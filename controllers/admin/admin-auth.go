package controllers

import (
	"os"
	"time"
	"errors"
	"strconv"
	"net/http"
	models "api/models"
	"github.com/rs/zerolog"
	token "api/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Define Logger Instance
var logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()

// login admin
func AdminLogin(c *gin.Context) {

	// Validate Input
	var input *models.AdminLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": err.Error(),
		})
		return
	}

	// Get Admin data of admin logging in
	var admin models.Admins
	result := models.DB.First(&admin, "username = ?", input.Username)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Invalid username or Password",
		})
		return
	}

	// Verify password given
	if err := models.VerifyPassword(admin.Password, input.Password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"message": "Invalid Password",
		})
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
	token, err := token.GenerateToken(TOKEN_EXPIRES_IN, admin.ID, os.Getenv("TOKEN_SECRET"))
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

// logout admin
func LogoutAdmin(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "You have been logged out!",
	})
}
