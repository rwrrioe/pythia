package rest_handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/domain/requests"
)

type AuthHandler struct {
	auth authn.SSOService
}

func NewAuthHandler(auth authn.SSOService) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

// /api/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req requests.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid body",
			"details": err.Error(),
		})
	}

	token, err := h.auth.Login(c, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, authn.ErrInvalidCredentials):
			c.JSON(401, gin.H{"error": "invalid credentials"})
		case errors.Is(err, authn.ErrUserAlreadyExists):
			c.JSON(409, gin.H{"error": "user already exists"})
		case errors.Is(err, authn.ErrSSOUnavailable):
			c.JSON(503, gin.H{"error": "sso unavailable"})
		default:
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(200, gin.H{"token": token})
}

// /api/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid body", "details": err.Error()})
		return
	}

	userID, err := h.auth.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, authn.ErrUserAlreadyExists):
			c.JSON(409, gin.H{"error": "user already exists"})
		case errors.Is(err, authn.ErrSSOUnavailable):
			c.JSON(503, gin.H{"error": "sso unavailable"})
		default:
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(201, gin.H{
		"user_id": userID,
	})

}
