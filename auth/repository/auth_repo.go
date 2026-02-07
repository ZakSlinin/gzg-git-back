package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ZakSlinin/gzg-git-back/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func CreateUser(db *sql.DB, ctx context.Context, username, email, password, fullname string) error {
	query := `INSERT INTO users (username, password, email, fullname) VALUES ($1, $2, $3, $4)`

	_, err := db.ExecContext(ctx, query, username, password, email, fullname)

	if err != nil {
		return err
	}

	return nil
}

func LoginUser(db *sql.DB, ctx context.Context, email, password string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var user model.User
	var passwordHash string

	err := db.QueryRowContext(ctx, query, email).Scan(&passwordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, errors.New("Invalid credentials")
	}

	return &user, nil
}
