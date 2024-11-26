package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"dot-gogat-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/marine-br/golib-logger/logger"
)

type ProxyHandler struct {
	consulService *services.ConsulService
}

func NewProxyHandler(consulService *services.ConsulService) *ProxyHandler {
	return &ProxyHandler{
		consulService: consulService,
	}
}

func (h *ProxyHandler) HandleProxy(c *gin.Context) {
	// Extrai o nome do serviço da URL
	serviceName := strings.TrimPrefix(c.Param("service"), "/")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service name is required"})
		return
	}

	// Descobre o serviço no Consul
	services, err := h.consulService.DiscoverService(serviceName)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to discover service %s: %v", serviceName, err))
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service discovery failed"})
		return
	}

	if len(services) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	// Seleciona o primeiro serviço disponível (pode ser implementado um load balancer aqui)
	service := services[0]
	serviceURL := fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)

	// Constrói a URL completa para o proxy
	path := c.Param("path")
	targetURL := fmt.Sprintf("%s%s", serviceURL, path)

	// Cria a requisição para o serviço
	proxyReq, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create proxy request: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create proxy request"})
		return
	}

	// Copia os headers originais
	copyHeaders(c.Request.Header, proxyReq.Header)

	// Realiza a requisição
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to proxy request: %v", err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to proxy request"})
		return
	}
	defer resp.Body.Close()

	// Copia os headers da resposta
	copyHeaders(resp.Header, c.Writer.Header())
	c.Status(resp.StatusCode)

	// Copia o corpo da resposta
	io.Copy(c.Writer, resp.Body)
}

func copyHeaders(src, dst http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
