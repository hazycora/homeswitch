package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.gay/h/homeswitch/internal/crypto"
	account_model "git.gay/h/homeswitch/internal/models/account"
	app_model "git.gay/h/homeswitch/internal/models/app"
	token_model "git.gay/h/homeswitch/internal/models/token"
	"git.gay/h/homeswitch/internal/router/mastoapi/form"

	"github.com/gin-gonic/gin"
	"github.com/hazycora/go-mcache"
	"github.com/rs/zerolog/log"
)

var codeCache = mcache.New()

type AuthorizeForm struct {
	ResponseType string `json:"response_type" form:"response_type" validate:"required,eq=code"`
	ClientID     string `json:"client_id" form:"client_id" validate:"required"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri" validate:"required"`
	Scope        string `json:"scope" form:"scope"`
	Lang         string `json:"lang" form:"lang"`
	// TODO: add support for force_login
}

func AuthorizeHandler(c *gin.Context) {
	var requestForm AuthorizeForm
	err := c.Bind(&requestForm)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = form.ValidateForm(requestForm)
	if err != nil {
		formError, ok := err.(form.FormError)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, formError)
		return
	}

	code, err := crypto.RandomString(32)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating code"))
		return
	}

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := token_model.GetToken(accessToken)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	codeCache.Set(code, *token.UserID, time.Hour)

	if requestForm.RedirectURI == "urn:ietf:wg:oauth:2.0:oob" {
		c.Writer.WriteString(fmt.Sprintf("Code: %s", code))
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s?code=%s", requestForm.RedirectURI, url.QueryEscape(code)))
}

type TokenForm struct {
	GrantType    string  `json:"grant_type" form:"grant_type" validate:"required"`
	Code         *string `json:"code" form:"code"`
	ClientID     string  `json:"client_id" form:"client_id" validate:"required"`
	ClientSecret string  `json:"client_secret" form:"client_secret" validate:"required"`
	RedirectURI  string  `json:"redirect_uri" form:"redirect_uri" validate:"required"`
	Scope        string  `json:"scope" form:"scope"`
}

func TokenHandler(c *gin.Context) {
	var requestForm TokenForm
	err := c.Bind(&requestForm)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = form.ValidateForm(requestForm)
	if err != nil {
		formError, ok := err.(form.FormError)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, formError)
		return
	}

	var userId *string
	if requestForm.Code != nil {
		code := *requestForm.Code
		cachedUserId, ok := codeCache.Get(code)
		if !ok {
			log.Debug().Msg("Invalid code")
			c.Status(http.StatusBadRequest)
			c.Writer.WriteString("Invalid code")
			return
		}
		userIdString := cachedUserId.(string)
		userId = &userIdString
		codeCache.Remove(code)
	}

	app, err := app_model.GetApp(requestForm.ClientID)
	if err != nil || app.ClientSecret != requestForm.ClientSecret {
		log.Error().Err(err).Str("client_id", requestForm.ClientID).Msg("Failed to get app")
		if err == nil {
			log.Debug().Str("secret", app.ClientSecret).Str("claimed_secret", requestForm.ClientSecret).Msg("Secrets")
		}
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := &token_model.Token{
		ClientID: requestForm.ClientID,
	}
	if userId != nil {
		_, ok := account_model.GetAccountByID(*userId)
		if !ok {
			log.Debug().Str("user_id", *userId).Msg("Account not found")
			c.Status(http.StatusBadRequest)
			c.Writer.WriteString("Account not found")
		}
		token.UserID = userId
	}
	err = token_model.CreateToken(token)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := map[string]interface{}{
		"access_token": token.AccessToken,
		"token_type":   "Bearer",
		"scope":        strings.Join(token.Scopes, " "),
		"created_at":   token.CreatedAt,
	}
	c.JSON(http.StatusOK, response)
}
