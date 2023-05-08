package main

import (
	"net/http"
	"api/models"
	"api/middlewares"
	"github.com/gin-gonic/gin"
	controllers "api/controllers"
)


func main(){

	// Initialzie DB Connection
	models.ConnectDB() 
	
	r := gin.Default() //router

	// Group endpoints
	public := r.Group("/api/auth")
	protected := r.Group("/api/users")

	// Test API works
	r.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Routes
	// users
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/me", controllers.GetCurrentUser)
	protected.GET("/me/points", controllers.GetPointsData)

	// Auth
	public.POST("/register", controllers.CreateAccount)
	public.GET("/verify-email/:secret_code", controllers.VerifyEmail)
	public.POST("/login", controllers.Login)

	r.Run() // listen and serve on 0.0.0.0:8080	
}
