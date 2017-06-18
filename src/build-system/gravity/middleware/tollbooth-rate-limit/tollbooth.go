package gintoll

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"github.com/gin-gonic/gin"
)

// LimitMiddleware wraps the tollbooth Limiter's functionality
// and aborts or calls next for subsequent handlers accordingly
func LimitMiddleware(limiter *config.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		tollbooth.SetResponseHeaders(limiter, c.Writer)

		if err := tollbooth.LimitByRequest(limiter, c.Request); err != nil {
			c.AbortWithError(err.StatusCode, err)
		}

		c.Next()
	}
}
