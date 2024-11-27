package routes

import (
	"dot-gogat-api/internal/handlers"
	"dot-gogat-api/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterProxyRoutes(r *gin.RouterGroup, consulService *services.ConsulService) {
	proxyHandler := handlers.NewProxyHandler(consulService)

	proxy := r.Group("/proxy")
	{
		proxy.Any("/:service/*path", proxyHandler.HandleProxy)
	}
}
