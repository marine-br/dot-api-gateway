package server

import (
	"dot-gogat-api/configs"
	"dot-gogat-api/internal/middlewares"
	"dot-gogat-api/internal/server/routes"
	"dot-gogat-api/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, consulService *services.ConsulService) {
	config := configs.LoadConfig()
	r.Use(middlewares.LoggerMiddleware())
	r.Use(middlewares.RateLimitMiddleware(config.RateLimit.Window, config.RateLimit.Limit, nil))
	routes.RegisterHealthRoutes(r)

	v1 := r.Group("/v1")
	{
		routes.RegisterAuthRoutes(v1)

		protected := v1.Group("")
		protected.Use(middlewares.AuthMiddleware(config.JWT.SecretKey))
		routes.RegisterProxyRoutes(protected, consulService)
	}
}
