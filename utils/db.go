package main

import(
	"fmt"
	"database/sql"
	"github.com/lib/pq"
	"github.com/joho/godotenv"
)

func ConnectDB(){

	// Load .env file
	godotenv.Load(".env")

	// Database connection credentials
	host := goDotEnvVariable("host")
    port := goDotEnvVariable("port")
    user := goDotEnvVariable("user")
    password := goDotEnvVariable("password")
    dbname := goDotEnvVariable("dbname")

	// connection string
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

	// Connect to Database
	db, err := sql.Open("postgres", connStr)
	if err != nil{
		panic(err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil{
		panic(err)
	}

	return db, nil
}
