package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	FullName         string    `json:"fullname"`
	Bio              *string   `json:"bio"`
	AvatarURL        *string   `json:"avatar_url"`
	PublicReposCount int       `json:"public_repos_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UserDTO struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	FullName         string    `json:"fullname"`
	Bio              *string   `json:"bio"`
	AvatarURL        *string   `json:"avatar_url"`
	PublicReposCount int       `json:"public_repos_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RegisterResponse struct {
	User        UserDTO `json:"user"`
	AccessToken string  `json:"access_token"`
}
