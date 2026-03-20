package middleware

import (
	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/model"
	"password-timer-manager/internal/pkg/response"
)

func RequirePermission(perm model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "not authenticated")
			c.Abort()
			return
		}

		u := user.(*model.User)
		if !u.HasPermission(perm) {
			response.Forbidden(c, "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
