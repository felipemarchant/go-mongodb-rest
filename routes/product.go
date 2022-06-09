package routes

import (
	"github.com/felipemarchant/go-mongo-rest/controllers"
	"github.com/gin-gonic/gin"
)

func Product(g *gin.Engine) {
	g.POST("/products", controllers.AddProduct)
	g.GET("/products", controllers.GetProducts)
	g.GET("/products/search", controllers.SearchProduct)
}
