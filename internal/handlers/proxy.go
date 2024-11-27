package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"

	"dot-gogat-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/marine-br/golib-logger/logger"
)

type ProxyHandler struct {
	consulService *services.ConsulService
	currentIndex  uint32
}

func NewProxyHandler(consulService *services.ConsulService) *ProxyHandler {
	return &ProxyHandler{
		consulService: consulService,
		currentIndex:  0,
	}
}

func (h *ProxyHandler) selectService(services []*api.ServiceEntry) *api.ServiceEntry {
	if len(services) == 0 {
		return nil
	}

	index := atomic.AddUint32(&h.currentIndex, 1) % uint32(len(services))
	return services[int(index)]
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

	// Seleciona por round-robin
	service := h.selectService(services)
	if service == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no healthy services available"})
		return
	}

	// Verifica se o serviço está saudável
	if service.Checks.AggregatedStatus() != api.HealthPassing {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service is not healthy"})
		return
	}

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
