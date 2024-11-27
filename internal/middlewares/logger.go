package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marine-br/golib-logger/logger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Extrai o userID do contexto (se existir)
		var userID string
		if id, exists := c.Get("userID"); exists {
			userID = fmt.Sprintf("%v", id)
		}

		// Processa a requisição
		c.Next()

		// Prepara os dados do log
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		// Formata a mensagem de log com base na presença do userID
		var logMessage string
		if userID != "" {
			logMessage = fmt.Sprintf("[%s] %s %s %d %v \"%s\" userID:%s",
				method,
				path,
				clientIP,
				statusCode,
				latency,
				userAgent,
				userID,
			)
		} else {
			logMessage = fmt.Sprintf("[%s] %s %s %d %v",
				method,
				path,
				clientIP,
				statusCode,
				latency,
			)
		}

		// Log com nível apropriado baseado no status code
		if statusCode >= 500 {
			logger.Error(logMessage)
		} else if statusCode >= 400 {
			logger.Warning(logMessage)
		} else {
			logger.Log(logMessage)
		}
	}
}
