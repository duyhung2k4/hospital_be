package config

import (
	"net/smtp"

	"gorm.io/gorm"
)

var (
	appPort string
	appHost string
	// pythonNodePort string

	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string

	smtpEmail    string
	smtpHost     string
	smtpPort     string
	smtpPassword string

	dbPsql *gorm.DB

	authSmtp smtp.Auth
)
