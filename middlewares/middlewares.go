package middleware

import (
	"api/models"
	tokens "api/utils/token"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func DeserializeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		cookie, err := c.Cookie("token")

		authorizationHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		// Load .env file
		errr := godotenv.Load(".env")
		if errr != nil {
			log.Fatal("Error loading .env file")
		}

		sub, err := tokens.ValidateToken(token, os.Getenv("TOKEN_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User
		result := models.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": "fail",
				"message": "the user belonging to this token no logger exists",
			})
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}

func DeserializeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		cookie, err := c.Cookie("token")

		authorizationHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"message": "You are not logged in",
			})
			return
		}

		// Load .env file
		errr := godotenv.Load(".env")
		if errr != nil {
			log.Fatal("Error loading .env file")
		}

		sub, err := tokens.ValidateToken(token, os.Getenv("TOKEN_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"message": err.Error(),
			})
			return
		}

		var admin models.Admins
		result := models.DB.First(&admin, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": "fail",
				"message": "Admin belonging to this token no longer exists",
			})
			return
		}

		c.Set("currentUser", admin)
		c.Next()
	}
}