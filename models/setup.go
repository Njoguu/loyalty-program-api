/*
In this setup.go, we use godotenv for fetching
environment variables
*/

package models

import (
	"os"
	"fmt"
	"log"
	"context"
	"strconv"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Define postgres Database
var DB *gorm.DB

// Define the Redis client
var RDB *redis.Client

// Define the cache instance
var MYCACHE *cache.Cache

// Initialize the Redis client
func InitRedisClient() {

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisDBStr := os.Getenv("REDISDATABASE")

	// convert port to integer
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISHOST"),
		Password: "", // Add password if required
		DB:       redisDB,  // Specify the Redis database index
	})

	// Ping the Redis server to check the connection
	_, errr := rdb.Ping(context.Background()).Result()
	if errr != nil{
		fmt.Printf("[REDIS-CLIENT]: Failed to connect to Redis: %v", err)
	}else{
		fmt.Println("[REDIS-CLIENT]: Redis connection initiated!")
	}

	RDB = rdb
}

// Initialize the cache instance
func InitCache() {
	mycache := cache.New(&cache.Options{
		Redis: RDB,
	})

	MYCACHE = mycache
}

// Initialize db instance
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
		fmt.Printf("[POSTGRES-DATABASE]: Cannot connect to database! %v", err)
		log.Fatal("Connection error:", err)
	}else{
		fmt.Println("[POSTGRES-DATABASE]: Database connection initiated!")
	}

	DB = db
}
