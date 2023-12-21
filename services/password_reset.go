package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mohieey/lenslocked/models"
	"github.com/mohieey/lenslocked/random"
)

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

func (prs *PasswordResetService) Create(email string) (*models.PasswordReset, error) {
	email = strings.ToLower(email)
	var userID int
	row := prs.DB.QueryRow(
		`
		SELECT id
		FROM users
		WHERE email = $1
		`, email,
	)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user from database %w", err)
	}

	token, err := random.SessionToken(int(math.Max(float64(prs.BytesPerToken), MinBytesPerToken)))
	if err != nil {
		return nil, fmt.Errorf("error creating password reset token: %w", err)
	}

	duration := int(math.Max(float64(prs.Duration), float64(models.MinResetDuration)))

	pwReset := models.PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: prs.hash(token),
		ExpiresAt: time.Now().Add(time.Duration(duration)),
	}
	row = prs.DB.QueryRow(
		`
		INSERT INTO password_resets(user_id, token_hash, expires_at)
		VALUES($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;
		`,
		pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt,
	)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("error inserting password reset: %w", err)
	}

	return &pwReset, nil
}

func (prs *PasswordResetService) Consume(token string) (*models.User, error) {
	tokenHash := prs.hash(token)
	var user models.User
	var pwReset models.PasswordReset

	row := prs.DB.QueryRow(`
			SELECT password_resets.id, password_resets.expires_at,
				users.id, users.email
			FROM password_resets
				JOIN users ON users.id = password_resets.user_id
			WHERE password_resets.token_hash = $1;`, tokenHash)
	err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	err = prs.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (prs *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (prs *PasswordResetService) delete(id int) error {
	_, err := prs.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("error deleting password reset token: %w", err)
	}
	return nil
}
