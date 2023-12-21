package services

import (
	"database/sql"
	"time"

	"github.com/mohieey/lenslocked/models"
)

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

func (prs *PasswordResetService) Create(email string) (*models.PasswordReset, error) {
	return nil, nil
}

func (prs *PasswordResetService) Consume(token string) (*models.User, error) {
	return nil, nil
}
