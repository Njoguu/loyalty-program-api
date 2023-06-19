package main

import (
	"net/http"
	"api/models"
	"api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	controllers "api/controllers"
	adminControllers "api/controllers/admin"
)


var Server *gin.Engine

func init(){
	// Initialzie DB Connection
	models.ConnectDB() 
	
	Server = gin.Default() //router
}	

func main(){
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true

	Server.Use(cors.New(corsConfig))

	// Group endpoints
	public := Server.Group("/api/auth")
	protected := Server.Group("/api/users")
	admins := Server.Group("/api/admin")
	admins_protected := Server.Group("/api/admin") 

	// Test API works
	Server.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Routes
	// users
	protected.Use(middleware.DeserializeUser())
	protected.GET("/me", controllers.GetCurrentUser)
	protected.POST("/me/redeem", controllers.RedeemPoints)
	protected.GET("/me/transaction-history", controllers.ViewTransactions)
	protected.GET("/products", controllers.GetProducts)
	protected.PATCH("/me/change-password", controllers.ChangePassword)

	// Auth
	public.POST("/register", controllers.CreateAccount)
	public.GET("/verify-email/:secret_code", controllers.VerifyEmail)
	public.POST("/login", controllers.Login)
	public.GET("/logout", controllers.Logout)
	public.GET("/sessions/oauth/google", controllers.GoogleOAuth)
	public.POST("/forgot-password", controllers.ForgotPassword) 
	public.PATCH("/reset-password/:resetToken", controllers.ResetPassword)

	// Admins
	admins_protected.Use(middleware.DeserializeAdmin())
	admins.POST("/login", adminControllers.AdminLogin)
	
	admins_protected.GET("/logout", adminControllers.LogoutAdmin)
	admins_protected.POST("/product", adminControllers.AddProduct)
	admins_protected.PUT("/product/:id", adminControllers.UpdateProduct)
	admins_protected.DELETE("/product/:id", adminControllers.DeleteProduct)
	admins_protected.GET("/product", adminControllers.GetAllProducts)

	Server.Run(":" + "8000") // listen and serve on 0.0.0.0:8000	
}
