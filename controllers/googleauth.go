package controllers

import (
	"os"
	"fmt"
	"log"
	"time"
	"strconv"
	"strings"
	"net/http"
	"api/utils"
	models "api/models"
	token "api/utils/token"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)


func GoogleOAuth(c *gin.Context) {
	code := c.Query("code")
	var pathUrl string = "/"

	if c.Query("state") != "" {
		pathUrl = c.Query("state")
	}

	if code == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"message": "Authorization code not provided!",
		})
		return
	}

	tokenRes, err := utils.GetGoogleOAuthToken(code)

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"status": "fail",
			"message": err.Error(),
		})
		return
	}

	google_user, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"status": "fail",
			"message": err.Error(),
		})
		return
	}

	email := strings.ToLower(google_user.EmailAddress)

	user_data := models.User{
		Username:      google_user.FirstName,
		EmailAddress:     email,
		FirstName: google_user.FirstName,
		LastName: google_user.LastName,
		Gender: google_user.Gender,
		Password:  "",
		City: google_user.City,
		PhoneNumber: google_user.PhoneNumber,
		IsEmailVerified: true,
	}

	if models.DB.Model(&user_data).Where("email = ?", email).Updates(&user_data).RowsAffected == 0 {
		models.DB.Create(&user_data)
	}

	var user models.User
	models.DB.First(&user, "email = ?", email)

	// Load .env file
	errr := godotenv.Load(".env")
	if errr != nil {
		log.Fatal("Error loading .env file")
	}

	TOKEN_EXPIRES_IN, err := time.ParseDuration(os.Getenv("TOKEN_EXPIRES_IN"))
	if err != nil{
		return 
	}

	JWTTokenSecret := os.Getenv("TOKEN_SECRET")

	TokenMaxage, err := strconv.Atoi(os.Getenv("TOKEN_MAXAGE"))
	if err != nil{
		return 
	}

	token, err := token.GenerateToken(TOKEN_EXPIRES_IN, user.ID, JWTTokenSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": err.Error(),
		})
		return
	}

	c.SetCookie("token", token, TokenMaxage*60, "/", "localhost", false, true)

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(os.Getenv("FRONTEND_ORIGIN"), pathUrl))
}

