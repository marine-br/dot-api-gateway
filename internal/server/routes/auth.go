package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", handleLogin)
		auth.POST("/refresh", handleRefreshToken)
		auth.GET("/verify", handleVerifyToken)
	}
}

func handleLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login endpoint - to be implemented",
	})
}

func handleRefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Refresh token endpoint - to be implemented",
	})
}

func handleVerifyToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Verify token endpoint - to be implemented",
	})
}
