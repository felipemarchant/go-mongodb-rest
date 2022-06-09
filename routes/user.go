package routes

import (
	"github.com/felipemarchant/go-mongo-rest/controllers"
	"github.com/gin-gonic/gin"
)

func User(g *gin.Engine) {
	g.POST("/users/signup", controllers.SignUp)
	g.POST("/users/signin", controllers.Login)
}
