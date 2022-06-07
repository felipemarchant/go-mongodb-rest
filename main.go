package main

import (
	"github.com/gin-gonic/gin"
	"go-mongo-rest/routes"
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
