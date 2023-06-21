package utils

import (
	"os"
	"fmt"
	"bytes"
	"errors"
	"strconv"
	"api/models"
	"crypto/tls"
	"html/template"
	"path/filepath"
	"gopkg.in/gomail.v2"
	"github.com/rs/zerolog"
	"github.com/k3a/html2text"
	"github.com/joho/godotenv"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// Define Logger Instance
var logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()

// ? Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.User, data *EmailData, templateName string) {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error().Err(errors.New("reading environment variables failed")).Msgf("%v",err)
	}

	// Sender data.
	from := os.Getenv("EMAIL_FROM")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpUser := os.Getenv("SMTP_USER")
	to := user.EmailAddress
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")

	// convert port to integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		logger.Error().Err(errors.New("converting port to integer failed")).Msgf("%v",err)
	}

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		logger.Error().Err(errors.New("could not parse template")).Msgf("%v",err)
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		logger.Error().Err(errors.New("sending email failed")).Msgf("%v",err)
	}
}
