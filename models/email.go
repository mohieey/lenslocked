package models

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
