package routes

import (
	"github.com/felipemarchant/go-mongo-rest/controllers"
	"github.com/gin-gonic/gin"
)

func Portal(g *gin.Engine) {
	HealthCheck(g)
	User(g)
	g.GET("/products", controllers.GetProducts)
	g.GET("/products/search", controllers.SearchProduct)
}
