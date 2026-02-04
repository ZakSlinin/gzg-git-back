package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

//TODO: написать реализацию после тестов

func (h *AuthHandler) Register(c *gin.Context) {
	c.JSON(201, gin.H{})
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(200, gin.H{})
}
