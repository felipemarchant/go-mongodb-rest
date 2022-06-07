package controllers

import (
	"github.com/felipemarchant/go-mongo-rest/database"
	"github.com/felipemarchant/go-mongo-rest/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()
	defer ctx.Done()

	response := gin.H{"status": "OK"}

	err := database.Client.Ping(ctx)
	if err != nil {
		response["status"] = "Internal Server Error"
		response["err"] = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusFound, response)
}
