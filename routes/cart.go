package routes

import (
	"github.com/felipemarchant/go-mongo-rest/controllers"
	"github.com/gin-gonic/gin"
)

func Cart(g *gin.Engine) {
	g.GET("/cart", controllers.ListCartItem)
	g.POST("/cart", controllers.AddToCart)
	g.DELETE("/cart/:product", controllers.DeleteFromCart)
	g.POST("/cart/checkout", controllers.CheckoutCart)
	g.POST("/cart/checkout/direct", controllers.InstantCheckout)
}
