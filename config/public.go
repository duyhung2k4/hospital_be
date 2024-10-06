package config

import (
	"net/smtp"

	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

func GetPsql() *gorm.DB {
	return dbPsql
}

func GetAppPort() string {
	return appPort
}

func GetAppHost() string {
	return appHost
}

func GetSmtpPort() string {
	return smtpPort
}

func GetSmtpHost() string {
	return smtpHost
}

func GetAuthSmtp() smtp.Auth {
	return authSmtp
}

func GetJWT() *jwtauth.JWTAuth {
	return jwt
}
