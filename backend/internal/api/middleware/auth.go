package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/model"
	"password-timer-manager/internal/pkg/response"
	"password-timer-manager/internal/service"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := m.authService.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		user := &model.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		c.Set("user", user)
		c.Next()
	}
}

func (m *AuthMiddleware) APIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			response.Unauthorized(c, "missing API key")
			c.Abort()
			return
		}

		user, err := m.authService.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			response.Unauthorized(c, "invalid API key")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := m.authService.ValidateToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		user := &model.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		c.Set("user", user)
		c.Next()
	}
}
