package main

import (
	"github.com/felipemarchant/go-mongo-rest/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8000"
	}

	router := gin.New()

	routes.HealthCheck(router)

	router.Use(gin.Logger())
	log.Fatal(router.Run(":" + port))
}
