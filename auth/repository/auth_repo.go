package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ZakSlinin/gzg-git-back/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

type AuthRepository interface {
	CreateUser(ctx context.Context, username, email, password, fullname string) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

func (repo *PostgresAuthRepository) CreateUser(db *sql.DB, ctx context.Context, username, email, password, fullname string) error {
	id := uuid.New()
	now := time.Now()

	query := `INSERT INTO users (id, username, email, password, fullname, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING id, username, email, fullname, bio, avatar_url, public_repos_count, created_at, updated_at`

	u := &model.User{}

	err := repo.db.QueryRowContext(ctx, query,
		id, username, email, password, fullname, now).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.FullName,
		&u.Bio,
		&u.AvatarURL,
		&u.PublicReposCount,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresAuthRepository) LoginUser(db *sql.DB, ctx context.Context, email, password string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var user model.User
	var passwordHash string

	err := repo.db.QueryRowContext(ctx, query, email).Scan(&passwordHash)

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

func (repo *PostgresAuthRepository) GetUserByEmail(db *sql.DB, ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, username, email, fullanme, bio, avatar_url, public_repos_count, created_at, updated_at  FROM users WHERE email = $1`

	u := &model.User{}

	row := repo.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.FullName,
		&u.Bio,
		&u.AvatarURL,
		&u.PublicReposCount,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return u, nil
}
