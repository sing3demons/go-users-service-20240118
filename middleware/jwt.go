package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/users/security"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Request.Header.Get("Authorization")
		if s == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}

		token := strings.TrimPrefix(s, "Bearer ")

		claims, err := security.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}
		sub, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}

		c.Set("userId", sub)
		c.Next()
	}
}
