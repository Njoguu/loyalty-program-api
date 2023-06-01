package main

import (
	"fmt"
	"log"
	"api/models"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect Golang server to the Postgres database.
	models.ConnectDB()
}

func main() {
	models.DB.AutoMigrate(&models.User{}, &models.VerifyEmails{},
	&models.Products{}, &models.Transactions{})
	fmt.Println("? Migration complete")
}
