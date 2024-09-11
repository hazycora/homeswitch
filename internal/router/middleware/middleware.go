package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AllowCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "accept, content-type, x-requested-with, authorization")
	c.Header("Access-Control-Allow-Methods", "POST, PUT, DELETE, GET, PATCH, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Link, X-RateLimit-Reset, X-RateLimit-Limit, X-RateLimit-Remaining, X-Request-Id")

	if c.Request.Method == http.MethodOptions && c.GetHeader("Access-Control-Request-Method") != "" {
		c.Status(http.StatusNoContent)
		c.Writer.WriteString("")
		return
	}

	c.Next()
}
