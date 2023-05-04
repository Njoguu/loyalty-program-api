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
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func ConnectDB() (*sql.DB, error){

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection credentials
	host := os.Getenv("host")
    portStr := os.Getenv("port")
    user := os.Getenv("user")
    password := os.Getenv("password")
    dbname := os.Getenv("dbname")

	// convert port to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	// connection string
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

	// Connect to Database
	db, err := sql.Open("postgres", connStr)
	if err != nil{
		panic(err)
	}

	// Test connection
	err = db.Ping()
	if err != nil{
		panic(err)
	}

	return db, nil
}
