package timelines

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// TODO: implement
	c.JSON(http.StatusOK, []string{})
}
