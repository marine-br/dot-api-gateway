package main

import (
	"fmt"
	"os"

	"dot-gogat-api/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/marine-br/golib-logger/logger"
)

func main() {
	fmt.Println("Api Gateway DoTelematics")

	r := gin.Default()
	r.SetTrustedProxies(nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server.RegisterRoutes(r)

	if err := r.Run(":" + port); err != nil {
		logger.Error("Erro ao iniciar o servidor: %v", err)
	}
}
