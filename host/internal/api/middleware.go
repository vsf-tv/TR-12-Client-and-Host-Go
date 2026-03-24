package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

// JWTAuth returns a gin middleware that validates JWT tokens.
func JWTAuth(accountSvc *service.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header", "code": 401})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := accountSvc.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token", "code": 401})
			c.Abort()
			return
		}
		c.Set("account_id", claims.AccountID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
