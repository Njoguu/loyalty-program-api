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
	public := r.Group("/api")
	protected := r.Group("/api/user")

	// Test API works
	r.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Public routes
	public.POST("/register", controllers.CreateAccount)
	public.POST("/login", controllers.Login)

	// Protected routes
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/me", controllers.GetCurrentUser)

	r.Run() // listen and serve on 0.0.0.0:8080	
}
