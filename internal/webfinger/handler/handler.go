package handler

import (
	"net/http"
	"regexp"

	"git.gay/h/homeswitch/internal/config"
	account_model "git.gay/h/homeswitch/internal/models/account"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	reResource = regexp.MustCompile(`^acct:([^@]+)@(.+)$`)
)

func Handler(c *gin.Context) {
	resource := c.Query("resource")
	if resource == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	submatches := reResource.FindStringSubmatch(resource)
	username := submatches[1]
	instance := submatches[2]
	if instance != config.ServerName {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	log.Debug().Str("username", username).Any("instance", instance).Msg("Webfinger request")
	account, ok := account_model.GetAccountByUsername(username)
	if !ok {
		http.Error(c.Writer, "Account not found", http.StatusNotFound)
		return
	}
	c.Header("Content-Type", "application/jrd+json")
	c.JSON(http.StatusOK, account.Webfinger())
}
