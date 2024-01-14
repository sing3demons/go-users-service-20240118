package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sing3demons/users/constant"
	"github.com/sing3demons/users/router"
	"github.com/sing3demons/users/security"
)

func Authorization() router.ServiceHandleFunc {
	return func(c router.IContext) {
		s := c.GetAuthorization()
		fmt.Println("GetAuthorization++++++++++++++++++>", s)
		if s == "" {
			// c.JSON(401, gin.H{"message": "unauthorized"})
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}

		has := strings.HasPrefix(s, constant.BEARER)
		if !has {
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}

		token := strings.TrimPrefix(s, constant.BEARER)

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

func Authorization2() gin.HandlerFunc {
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
