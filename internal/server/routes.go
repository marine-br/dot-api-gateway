package server

import (
	"dot-gogat-api/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "API Gateway is up and running",
		})
	})

	r.Use(middlewares.ProxyMiddleware())
}
