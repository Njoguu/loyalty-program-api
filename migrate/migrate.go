package main

import (
	"os"
	"errors"
	"api/models"
	"github.com/rs/zerolog"
	"github.com/joho/godotenv"
)

// Define Logger Instance
var logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()

func init() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error().Err(errors.New("reading environment variables failed")).Msgf("%v",err)
	}

	// connect Golang server to the Postgres database.
	models.ConnectDB()
}

func main() {

	logger.Info().Msg("starting db migration")

	models.DB.AutoMigrate(&models.User{}, &models.VerifyEmails{},
	&models.Products{}, &models.Transactions{}, &models.Admins{},
	&models.PasswordReset{})

	logger.Info().Msg("db migration complete")
}
