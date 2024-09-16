package apicontext

import (
	"github.com/gin-gonic/gin"

	account_model "git.gay/h/homeswitch/internal/models/account"
	app_model "git.gay/h/homeswitch/internal/models/app"
)

func GetAccount(c *gin.Context) (account *account_model.Account, ok bool) {
	accountAny, ok := c.Get(AccountContextKey)
	if !ok || accountAny == nil {
		ok = false
		return
	}
	account = accountAny.(*account_model.Account)
	ok = true
	return
}

func GetApp(c *gin.Context) (app *app_model.App, ok bool) {
	appAny, ok := c.Get(AppContextKey)
	if !ok || appAny == nil {
		ok = false
		return
	}
	app = appAny.(*app_model.App)
	ok = true
	return
}
