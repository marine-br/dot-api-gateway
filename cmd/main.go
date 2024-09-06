package main

import (
	"fmt"
	"log"
	"os"

	"dot-gogat-api/internal/server"

	"github.com/gin-gonic/gin"
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
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
