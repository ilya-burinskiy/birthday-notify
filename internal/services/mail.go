package services

import (
	"fmt"
	"net/smtp"

	"github.com/ilya-burinskiy/birthday-notify/internal/configs"
)

type EmailSender struct {
	host     string
	port     string
	userName string
	auth     smtp.Auth
}

func NewEmailSender(config configs.Config) EmailSender {
	return EmailSender{
		host:     config.SMTPHost,
		port:     config.SMTPPort,
		userName: config.SMTPAuthUsername,
		auth: smtp.PlainAuth(
			config.SMTPAuthIdentity,
			config.SMTPAuthUsername,
			config.SMTPAuthPassword,
			config.SMTPHost,
		),
	}
}

func (sender EmailSender) Send(to string, subject string, body string) error {
	smtpAddr := sender.host + ":" + sender.port
	msg := subject + "\n" + body
	if err := smtp.SendMail(smtpAddr, sender.auth, sender.userName, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
