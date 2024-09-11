package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"git.gay/h/homeswitch/internal/config"
	account_model "git.gay/h/homeswitch/internal/models/account"
	app_model "git.gay/h/homeswitch/internal/models/app"
	token_model "git.gay/h/homeswitch/internal/models/token"
)

type SignInForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func SignInHandler(c *gin.Context) {
	var form SignInForm
	err := c.Bind(&form)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, ok := account_model.AccountLogin(form.Email, form.Password)
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}
	token := &token_model.Token{
		ClientID:  app_model.SystemApp.ClientID,
		TokenType: "Bearer",
		UserID:    &user.ID,
	}
	err = token_model.CreateToken(token)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetCookie("access_token", token.AccessToken, int(time.Hour.Seconds()*24*30), "/", config.ServerName, true, true)
	c.Redirect(http.StatusSeeOther, "/")
}
