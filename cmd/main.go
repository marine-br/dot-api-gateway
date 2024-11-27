package main

import (
	"dot-gogat-api/configs"
	"dot-gogat-api/internal/server"
	"dot-gogat-api/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/marine-br/golib-logger/logger"
)

func main() {
	config := configs.LoadConfig()
	consulService, err := services.NewConsulService(config.Consul.Address)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create Consul client: %v", err))
		return
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to register service in Consul: %v", err))
		return
	}

	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Success("API Gateway Starting...")
	r := gin.New()
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	server.RegisterRoutes(r, consulService)
	logger.Log(fmt.Sprintf("Server starting on port: %s", config.Port))

	if err := r.Run(":" + config.Port); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
