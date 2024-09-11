package activitypub

import (
	"fmt"
	"net/http"

	account_model "git.gay/h/homeswitch/internal/models/account"
	"github.com/gin-gonic/gin"
)

func ActorHandler(c *gin.Context) {
	username := c.Param("username")
	account, ok := account_model.GetAccountByUsername(username)
	if !ok {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("account not found"))
		return
	}
	c.Header("Content-Type", "application/activity+json")
	response := account.ActivityPub()
	c.JSON(http.StatusOK, response)
}
