package rest

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, payload any, status int) {
	c.JSON(status, gin.H{"status": status, "data": payload})
}
