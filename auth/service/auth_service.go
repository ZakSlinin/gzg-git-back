package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ZakSlinin/gzg-git-back/model"
	"github.com/ZakSlinin/gzg-git-back/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type JWTManager struct {
	secret     string
	expireTime time.Duration
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTManager() (*JWTManager, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrNoJwtSecret
	}

	expireTimeStr := os.Getenv("JWT_EXPIRE_TIME")
	if expireTimeStr == "" {
		expireTimeStr = "24h" // default
	}

	expireTime, err := time.ParseDuration(expireTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRE_TIME format: %w", err)
	}

	return &JWTManager{secret, expireTime}, err
}

func (m *JWTManager) GenerateToken(userID, username, email string) (string, error) {
	now := time.Now()
	exp := now.Add(m.expireTime)
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp), // When expires
			IssuedAt:  jwt.NewNumericDate(now), // When created
			NotBefore: jwt.NewNumericDate(now), // When start
			Issuer:    "gzg-git-back",          // Who create
			Subject:   userID,                  // For who
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var (
	ErrEmailAlreadyExist    = errors.New("email is already exists")
	ErrNoJwtSecret          = errors.New("jwt secret not set")
	ErrUsernameAlreadyExist = errors.New("username is already exists")
)

type AuthService struct {
	authUser repository.AuthRepository
}

func NewAuthService(authUser repository.AuthRepository) *AuthService {
	return &AuthService{authUser: authUser}
}

func (authService *AuthService) CreateUser(ctx context.Context, username, email, password, fullname string, bio, avatarUrl, createdAt *string) (*model.RegisterResponse, error) {
	receivedUser, err := authService.authUser.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if receivedUser != nil {
		if receivedUser.Username == username {
			return nil, ErrUsernameAlreadyExist
		}
		return nil, ErrEmailAlreadyExist
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	createdUser, err := authService.authUser.CreateUser(ctx, username, email, string(hash), fullname)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userDTO := &model.UserDTO{
		ID:               createdUser.ID,
		Username:         createdUser.Username,
		Email:            createdUser.Email,
		FullName:         createdUser.FullName,
		Bio:              createdUser.Bio,
		AvatarURL:        createdUser.AvatarURL,
		PublicReposCount: createdUser.PublicReposCount,
		CreatedAt:        createdUser.CreatedAt,
		UpdatedAt:        createdUser.UpdatedAt,
	}

	jwtManager, err := NewJWTManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT manager: %w", err)
	}

	token, err := jwtManager.GenerateToken(createdUser.ID.String(), username, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	resp := &model.RegisterResponse{
		User:        *userDTO,
		AccessToken: token,
	}

	return resp, nil
}
