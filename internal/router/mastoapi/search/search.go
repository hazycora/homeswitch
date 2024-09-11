package search

import (
	"net/http"
	"strings"

	account_model "git.gay/h/homeswitch/internal/models/account"
	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.Status(http.StatusBadRequest)
		c.Writer.WriteString("param is missing or the value is empty")
		return
	}

	query = strings.TrimSpace(query)

	if account_model.AcctRegExp.MatchString(query) {
		account_model.GetAccountByUsername(query)
	}

	// TODO: implement
}
