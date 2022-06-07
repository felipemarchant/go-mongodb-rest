package routes

import (
	"github.com/felipemarchant/go-mongo-rest/controllers"
	"github.com/gin-gonic/gin"
)

func Address(g *gin.Engine) {
	g.POST("/addresses", controllers.AddAddress)
	g.PUT("/addresses/homeaddress", controllers.EditHomeAddress)
	g.PUT("/addresses/workaddress", controllers.EditWorkAddress)
	g.DELETE("/addresses", controllers.DeleteAddress)
}
