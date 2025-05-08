package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func RegisterPingHandler(router *gin.Engine, logger *slog.Logger) {
	router.GET("/ping", pingHandler)
}
