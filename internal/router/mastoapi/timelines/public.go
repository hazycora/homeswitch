package timelines

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublicHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// TODO: implement
	c.JSON(http.StatusOK, []string{})
}
