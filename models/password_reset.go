package models

import (
	"time"
)

const (
	MinResetDuration = 10 * time.Minute
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}
