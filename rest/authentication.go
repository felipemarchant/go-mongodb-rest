package rest

import (
	"github.com/felipemarchant/go-mongo-rest/security"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			Response(c, "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims, err := security.ValidateToken(token)
		if err != "" {
			Response(c, "Forbidden", http.StatusForbidden)
			c.Abort()
			return
		}

		c.Set("UserPrincipal", claims)
		c.Next()
	}
}
