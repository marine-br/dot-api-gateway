package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "API Gateway is up and running",
			"time":   time.Now(),
		})
	})
}
