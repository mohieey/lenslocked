package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/mohieey/lenslocked/models"
	"github.com/mohieey/lenslocked/random"
)

const MinBytesPerToken = 32

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*models.Session, error) {
	token, err := random.SessionToken(int(math.Max(float64(ss.BytesPerToken), MinBytesPerToken)))
	if err != nil {
		return nil, fmt.Errorf("error creating session token: %w", err)
	}

	session := models.Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	row := ss.DB.QueryRow(
		`
		INSERT INTO sessions(user_id, token_hash)
		VALUES($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id;
		`,
		session.UserID, session.TokenHash,
	)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*models.User, error) {
	tokenHash := ss.hash(token)

	var user models.User
	row := ss.DB.QueryRow(
		`
		SELECT users.id, users.email
		From sessions
    		Join users on sessions.user_id = users.id
		WHERE sessions.token_hash = $1
		`,
		tokenHash,
	)
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(
		`
		DELETE FROM sessions WHERE token_hash = $1
		`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
