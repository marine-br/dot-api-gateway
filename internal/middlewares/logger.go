package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/marine-br/golib-logger/logger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Executa os handlers primeiro
		c.Next()

		// Extrai o userID do contexto (se existir)
		var userID string
		if id, exists := c.Get("userID"); exists {
			userID = fmt.Sprintf("%v", id)
		}

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		// Formata a mensagem de log com base na presença do userID
		var logMessage string
		if userID != "" {
			logMessage = fmt.Sprintf("[%s] %s %s %v \"%s\" userID:%s",
				method,
				path,
				clientIP,
				statusCode,
				userAgent,
				userID,
			)
		} else {
			logMessage = fmt.Sprintf("[%s] %s %s %v \"%s\"", method, path, clientIP, statusCode, userAgent)
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
