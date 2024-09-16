package instance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrendingTagsHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// TODO: implement
	c.JSON(http.StatusOK, []string{})
}

func TrendingStatusesHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// TODO: implement
	c.JSON(http.StatusOK, []string{})
}
