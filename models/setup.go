/*
In this setup.go, we use godotenv for fetching
environment variables
*/

package models

import (
	"os"
	"fmt"
	"errors"
	"context"
	"strconv"
	"github.com/rs/zerolog"
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

// Define Logger Instance
var logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()

// Initialize the Redis client
func InitRedisClient() {

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error().Err(errors.New("reading environment variables failed")).Msgf("%v",err)
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
		logger.Error().Err(errors.New("connection to redis instance failed")).Msgf("%v",errr)
	}else{
		logger.Info().Msg("redis connection initiated")
	}

	RDB = rdb
}

// Initialize the cache instance
func InitCache() {
	mycache := cache.New(&cache.Options{
		Redis: RDB,
	})

	logger.Info().Msg("cache instance initiated")

	MYCACHE = mycache
}

// Initialize db instance
func ConnectDB(){

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error().Err(errors.New("reading environment variables failed")).Msgf("%v",err)
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
		logger.Error().Err(errors.New("converting port to integer failed")).Msgf("%v",err)
	}

	// connection string
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

	// Connect to Database
	db, err := gorm.Open("postgres", connStr)
	if err != nil{
		logger.Error().Err(errors.New("connection to postgres database failed")).Msgf("%v",err)
	}else{
		logger.Info().Msg("postgres database connection initiated")
	}

	DB = db
}
