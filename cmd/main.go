package main

import (
	"fmt"
	"os"

	"dot-gogat-api/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/marine-br/golib-logger/logger"
)

func main() {
	logger.Success("Api Gateway DoTelematics")

	r := gin.Default()
	r.SetTrustedProxies(nil)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
		logger.Log(fmt.Sprintf("Listenig port: %v", port))
	}

	server.RegisterRoutes(r)

	if err := r.Run(":" + port); err != nil {
		logger.Error(fmt.Sprintf("Erro ao iniciar o servidor: %d", err))
	}
}
