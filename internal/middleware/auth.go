package middleware

import (
	"net/http"
	"strings"

	"auth/internal/services/tokenizer"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckJWT(tokenizerSrv tokenizer.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing token"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token format"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := tokenizerSrv.ParseClaims(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token"})
			c.Abort()
			return
		}

		userIdStr, ok := claims["x-user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid userId in token"})
			c.Abort()
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid userID format"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
