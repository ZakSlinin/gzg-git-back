package handler

import (
	"github.com/ZakSlinin/gzg-git-back/internal/auth/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"fullname"`
	AvatarUrl string `json:"avatar_url"`
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{400, "invalid request body"})
		return
	}

	avatarUrl, err := h.UploadAvatar(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{500, "failed to upload avatar"})
		return
	}

	resp, err := h.authService.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password, req.FullName, avatarUrl)
	if err != nil {
		if err == service.ErrEmailAlreadyExist {
			c.JSON(http.StatusBadRequest, ErrorResponse{400, "email already exists"})
			return
		}
		if err == service.ErrUsernameAlreadyExist {
			c.JSON(http.StatusBadRequest, ErrorResponse{400, "username already exists"})
			return
		}
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{400, "invalid request body"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{http.StatusUnauthorized, "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) UploadAvatar(c *gin.Context) (string, error) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{400, "invalid request file"})
		return "", err
	}

	if file.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, ErrorResponse{400, "file too large"})
		return "", err
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)

	err = c.SaveUploadedFile(file, "uploads/avatar/"+filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{400, "failed to save file"})
		return "", err
	}

	c.JSON(http.StatusOK, gin.H{"filename": filename})
	return filename, nil
}
