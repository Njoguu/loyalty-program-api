package main

import (
	"net/http"
	"api/models"
	"github.com/gin-gonic/gin"
	controllers "api/controllers"
)


func main(){
	r := gin.Default() //router

	db,err := models.ConnectDB() 
	if err != nil{
		panic(err)
	}
	// Provide db variable to controllers
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Test purposes
	r.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/user", controllers.AddUser)

	r.Run() // listen and serve on 0.0.0.0:8080	
}
