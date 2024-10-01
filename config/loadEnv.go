package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appHost = os.Getenv(APP_HOST)
	appPort = os.Getenv(APP_PORT)

	dbHost = os.Getenv(DB_HOST)
	dbPort = os.Getenv(DB_PORT)
	dbName = os.Getenv(DB_NAME)
	dbUser = os.Getenv(DB_USER)
	dbPassword = os.Getenv(DB_PASSWORD)

	smtpEmail = os.Getenv(SMTP_EMAIL)
	smtpHost = os.Getenv(SMTP_HOST)
	smtpPort = os.Getenv(SMTP_PORT)
	smtpPassword = os.Getenv(SMTP_PASSWORD)
}
