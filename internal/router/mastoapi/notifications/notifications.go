package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// TODO: implement
	c.JSON(http.StatusOK, []string{})
}
