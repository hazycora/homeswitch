package instance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomEmojiHandler(c *gin.Context) {
	emoji := []interface{}{}
	c.JSON(http.StatusOK, emoji)
}
