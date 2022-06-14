package security

import "github.com/gin-gonic/gin"

func UserPrincipal(c *gin.Context) *SignedDetails {
	if userPrinciapal, ok := c.Get("UserPrincipal"); ok {
		return userPrinciapal.(*SignedDetails)
	}
	return nil
}
