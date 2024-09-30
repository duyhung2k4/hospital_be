package utils

import (
	"app/config"
	"net/smtp"
)

type smtpUtils struct {
	smpt     smtp.Auth
	authSmtp smtp.Auth
	smtpHost string
	smtpPort string
}

type SmtpUtils interface {
	SendEmail(data string, email string) error
}

func (u *smtpUtils) SendEmail(data string, email string) error {
	to := []string{email}
	msg := []byte(data)

	err := smtp.SendMail(u.smtpHost+":"+u.smtpPort, u.authSmtp, email, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func NewSmtpUtils() SmtpUtils {
	return &smtpUtils{
		smpt:     config.GetAuthSmtp(),
		authSmtp: config.GetAuthSmtp(),
		smtpHost: config.GetSmtpHost(),
		smtpPort: config.GetSmtpPort(),
	}
}
