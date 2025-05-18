package adaptlimit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinMiddleware returns a Gin middleware function that uses the limiter
func (l *Limiter) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		release, err := l.Acquire(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		defer func() {
			// 如果状态码小于 500，我们认为请求是成功的
			success := c.Writer.Status() < 500
			release(success)
		}()

		c.Next()
	}
}
