package main

import (
	"net/http"
	"api/models"
	"api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	controllers "api/controllers"
)


var server *gin.Engine

func init(){
	// Initialzie DB Connection
	models.ConnectDB() 
	
	server = gin.Default() //router
}	

func main(){
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	// Group endpoints
	public := server.Group("/api/auth")
	protected := server.Group("/api/users")

	// Test API works
	server.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Routes
	// users
	protected.Use(middleware.DeserializeUser())
	protected.GET("/me", controllers.GetCurrentUser)

	// Auth
	public.POST("/register", controllers.CreateAccount)
	public.GET("/verify-email/:secret_code", controllers.VerifyEmail)
	public.POST("/login", controllers.Login)
	public.GET("/logout", controllers.Logout)
	public.GET("/sessions/oauth/google", controllers.GoogleOAuth)

	server.Run(":" + "8000") // listen and serve on 0.0.0.0:8000	
}
