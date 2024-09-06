package middlewares

import (
	"dot-gogat-api/internal/proxy"
	"io"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

func ProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		backendURL, found := proxy.FindBackend(c.Request.Method, c.Request.URL.Path)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			c.Abort()
			return
		}

		proxyReq, err := http.NewRequest(c.Request.Method, backendURL, c.Request.Body)
		if err != nil {
			log.Printf("Failed to create proxy request: %w", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to proxy request"})
			c.Abort()
			return
		}

		for key, values := range c.Request.Header {
			for _, v := range values {
				proxyReq.Header.Add(key, v)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			log.Printf("Failed to forward request to backend: %v", err)
			c.JSON(http.StatusBadGateway, gin.H{"err": "failed to forward request"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		c.Status(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}
		io.Copy(c.Writer, resp.Body)
	}
}
