package routes

import (
	"github.com/felipemarchant/go-mongo-rest/controllers"
	"github.com/gin-gonic/gin"
)

func HealthCheck(g *gin.Engine) {
	g.GET("/health", controllers.HealthCheck)
}
