package controllers

import (
	"github.com/felipemarchant/go-mongo-rest/database"
	"github.com/felipemarchant/go-mongo-rest/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(c *gin.Context) {
	ctx, cancel := utils.PingContextWithTimeout()
	defer cancel()
	defer ctx.Done()

	response := gin.H{"status": http.StatusOK, "data": "OK"}

	err := database.Client.Ping(ctx)
	if err != nil {
		response["status"] = http.StatusInternalServerError
		response["data"] = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
