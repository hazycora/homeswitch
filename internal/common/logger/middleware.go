package logger

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Middleware(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil && r != http.ErrAbortHandler {
			log.Error().Interface("recover", r).Msg("requestPanic")
			fmt.Println(string(debug.Stack()))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		log.Debug().Fields(map[string]interface{}{
			"remoteAddr": c.Request.RemoteAddr,
			"path":       c.Request.URL.Path,
			"proto":      c.Request.Proto,
			"method":     c.Request.Method,
			"userAgent":  c.Request.UserAgent(),
			"statusCode": c.Writer.Status(),
			"bytesIn":    c.Request.ContentLength,
			"bytesOut":   c.Writer.Size(),
		}).Msg("request")
	}()
	c.Next()
}
