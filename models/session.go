package models

const MinBytesPerToken = 32

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}
