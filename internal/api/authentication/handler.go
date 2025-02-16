package authentication

import (
	"context"
	"net/http"

	"auth/internal/services/authentication"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthHandler interface {
	Authenticate(c *gin.Context)
}

type handler struct {
	authenticationSrv authentication.Service
}

func New(
	authenticationSrv authentication.Service,
) AuthHandler {
	return &handler{
		authenticationSrv: authenticationSrv,
	}
}

func (h *handler) Authenticate(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token, err := h.authenticationSrv.Authenticate(context.Background(), req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
