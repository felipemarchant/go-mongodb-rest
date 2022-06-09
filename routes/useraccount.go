package routes

import "github.com/gin-gonic/gin"

func UserAccount(g *gin.Engine) {
	Address(g)
	Product(g)
}
