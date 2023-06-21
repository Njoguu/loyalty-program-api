/*
A middleware function that logs all incoming HTTP requests to the server.
It adds contextual information like the request data, response codes, and more to each log entry.
*/

package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    log "github.com/rs/zerolog/log"
)


func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

		// Process the request
		c.Next()

		// Logging
		end := time.Now()
		latency := end.Sub(start)

        // Log the request details
		logger := log.With().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.Path).
			Str("user_agent", c.Request.UserAgent()).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Logger()

		logger.Info().Msg("Request handled")
    }
}

