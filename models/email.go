package models

import (
	"fmt"

	"github.com/go-mail/mail"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

type SMPTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMPTPConfig) *EmailService {
	return &EmailService{
		DefaultSender: "",
		dailer:        mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

}

type EmailService struct {
	DefaultSender string

	dailer *mail.Dialer
}

func (es *EmailService) Send(email Email) error {
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

func (es *EmailService) from(email *Email) string {
	switch {
	case email.From != "":
		return email.From
	case es.DefaultSender != "":
		return es.DefaultSender
	default:
		return DefaultSender
	}
}

func (es *EmailService) ForgotPassword(to, OTP string) error {
	email := Email{
		To:        to,
		Subject:   "Reset password",
		PlainText: "Use this OTP to reset your password",
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("error sending OTP %w", err)
	}

	return nil
}
