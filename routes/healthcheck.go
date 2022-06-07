package routes

import (
	"github.com/gin-gonic/gin"
	"go-mongo-rest/controllers"
)

func HealthCheck(g *gin.Engine) {
	g.GET("/health", controllers.HealthCheck)
}
