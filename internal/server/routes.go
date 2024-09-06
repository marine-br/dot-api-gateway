package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra as rotas da API Gateway
func RegisterRoutes(r *gin.Engine) {
	// Rota de exemplo para verificar o funcionamento da API Gateway
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "API Gateway is up and running",
		})
	})

	// Aqui você pode definir as rotas de proxy para diferentes microserviços
}
