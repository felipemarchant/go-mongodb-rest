package controllers

import (
	"fmt"
	"github.com/felipemarchant/go-mongo-rest/database"
	"github.com/felipemarchant/go-mongo-rest/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	response := map[string]interface{}{"status": "OK"}

	err := database.Client.Ping(ctx)
	if err != nil {
		response["status"] = "Internal Server Error"
		response["err"] = fmt.Sprintf("%s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusFound, response)
}
