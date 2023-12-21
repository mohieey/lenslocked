package services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mohieey/lenslocked/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*models.User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	row := us.DB.QueryRow(
		`
		INSERT INTO users(email, password_hash)
		VALUES($1, $2) RETURNING id;
		`,
		user.Email, user.PasswordHash,
	)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*models.User, error) {
	user := models.User{Email: strings.ToLower(email)}

	row := us.DB.QueryRow(
		`
		SELECT id, password_hash from users where email = $1;
		`,
		user.Email,
	)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password %w", err)

	}

	return &user, nil
}
