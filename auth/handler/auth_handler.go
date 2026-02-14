package handler

import (
	"github.com/ZakSlinin/gzg-git-back/service"
	"github.com/gin-gonic/gin"
	"net/http"
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

	resp, err := h.authService.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password, req.FullName, req.AvatarUrl)
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
