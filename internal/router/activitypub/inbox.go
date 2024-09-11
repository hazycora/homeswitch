package activitypub

import (
	"io"
	"net/http"

	account_model "git.gay/h/homeswitch/internal/models/account"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func InboxHandler(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	username := c.Param("username")
	account, ok := account_model.GetAccountByUsername(username)
	if !ok {
		http.Error(c.Writer, "Account not found", http.StatusNotFound)
		return
	}
	log.Debug().Str("body", string(body)).Str("for_account", account.Username).Msg("Inbox event, got body")
	// TODO: actually handle inbox events
}
