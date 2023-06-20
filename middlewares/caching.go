/*
This file contains the implementation of a caching middleware for API endpoints.
The caching middleware is responsible for caching the responses of certain routes 
to improve performance and reduce unnecessary database or computation overhead. 
It utilizes the go-redis library for efficient caching using Redis as the cache store. 
By integrating this caching middleware into the API, you can significantly enhance the 
overall responsiveness and scalability of the application.
*/

package middleware

import (
    "time"
	"bytes"
	"net/http"
	"api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
)

// Middleware function to cache API responses
func CacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique cache key based on the request URL
		cacheKey := c.Request.URL.String()

		// Check if the response is already cached
        var resp []byte
		err := models.MYCACHE.Get(c.Request.Context(), cacheKey, &resp)
		if err == nil {
			// If cached, return the cached response and skip the handler
			c.Data(http.StatusOK, "application/json", resp)
			c.Abort()
			return
		}

		// Create a custom response writer to capture the response body
		w := NewResponseWriter(c.Writer)

		// Replace the original response writer with the custom response writer
		c.Writer = w

		// Process the request by calling the next handler
		c.Next()

		// Cache the response if the status code is 200 (OK)
		if c.Writer.Status() == http.StatusOK {
			// Store the response body in the cache
			err := models.MYCACHE.Set(&cache.Item{
				Ctx:   c.Request.Context(),
				Key:   cacheKey,
				Value: w.Bytes(),
				TTL:   time.Hour, // Set the cache TTL as needed
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

// Custom response writer to capture the response body
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Create a new instance of the custom response writer
func NewResponseWriter(w gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
	}
}

// Override the Write method to capture the response body
func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// Bytes returns the captured response body
func (w *ResponseWriter) Bytes() []byte {
	return w.body.Bytes()
}
