package services

import (
	"fmt"

	"github.com/mohieey/lenslocked/models"
	"gopkg.in/mail.v2"
)

func NewEmailService(config models.SMPTPConfig) *EmailService {
	return &EmailService{
		DefaultSender: "",
		dailer:        mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

}

type EmailService struct {
	DefaultSender string

	dailer *mail.Dialer
}

func (es *EmailService) Send(email models.Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", es.from(&email))
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.PlainText != "" && email.HTML != "":
		msg.SetBody("text/plain", email.PlainText)
		msg.AddAlternative("text/html", email.HTML)
	case email.PlainText != "":
		msg.SetBody("text/plain", email.PlainText)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}
	msg.SetBody("text/plain", email.PlainText)
	msg.AddAlternative("text/html", email.HTML)

	err := es.dailer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

func (es *EmailService) from(email *models.Email) string {
	switch {
	case email.From != "":
		return email.From
	case es.DefaultSender != "":
		return es.DefaultSender
	default:
		return models.DefaultSender
	}
}

func (es *EmailService) ForgotPassword(to, resetUrl string) error {
	email := models.Email{
		To:        to,
		Subject:   "Reset password",
		PlainText: "Use this url to reset your password",
		HTML:      fmt.Sprintf(`<a href="%v">Reset Password</a>`, resetUrl),
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("error sending reset url %w", err)
	}

	return nil
}
