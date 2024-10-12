package handlers

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Ping example
// @Description Do ping
// @Tags User
// @Produce json
// @Success 200 {string} pong
// @Security AccessToken
// @Router /user/ping [get]
func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "pong")
	}
}
