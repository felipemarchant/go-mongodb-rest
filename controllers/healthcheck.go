package controllers

import (
	"github.com/felipemarchant/go-mongo-rest/database"
	r "github.com/felipemarchant/go-mongo-rest/rest"
	"github.com/felipemarchant/go-mongo-rest/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(c *gin.Context) {
	ctx, cancel := utils.PingContextWithTimeout()
	defer cancel()
	defer ctx.Done()
	if err := database.Client.Ping(ctx); err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Response(c, "OK", http.StatusOK)
}
