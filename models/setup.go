/*
In this setup.go, we use godotenv for fetching
environment variables
*/

package models

import (
	"os"
	"fmt"
	"log"
	"strconv"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_"github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func ConnectDB(){

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Database connection credentials
	host := os.Getenv("PG_HOST")
    portStr := os.Getenv("PG_PORT")
    user := os.Getenv("PG_USER")
    password := os.Getenv("PG_PASSWORD")
    dbname := os.Getenv("PG_DB")

	// convert port to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return
	}

	// connection string
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

	// Connect to Database
	db, err := gorm.Open("postgres", connStr)
	if err != nil{
		fmt.Println("Cannot connect to database!")
		log.Fatal("Connection error:", err)
	}else{
		fmt.Println("Database connection initiated!")
	}

	DB = db

	DB.AutoMigrate(&User{})
}
