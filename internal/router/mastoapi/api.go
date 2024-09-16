package mastoapi

import (
	"net/http"
	"strings"

	account_model "git.gay/h/homeswitch/internal/models/account"
	app_model "git.gay/h/homeswitch/internal/models/app"
	token_model "git.gay/h/homeswitch/internal/models/token"
	"git.gay/h/homeswitch/internal/router/mastoapi/accounts"
	"git.gay/h/homeswitch/internal/router/mastoapi/apicontext"
	"git.gay/h/homeswitch/internal/router/mastoapi/apps"
	"git.gay/h/homeswitch/internal/router/mastoapi/instance"
	instance_v1 "git.gay/h/homeswitch/internal/router/mastoapi/instance/v1"
	instance_v2 "git.gay/h/homeswitch/internal/router/mastoapi/instance/v2"
	"git.gay/h/homeswitch/internal/router/mastoapi/notifications"
	"git.gay/h/homeswitch/internal/router/mastoapi/timelines"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RequirePostJSONBody(c *gin.Context) {
	if c.Request.Method == "POST" && c.GetHeader("Content-Type") != "application/json" {
		http.Error(c.Writer, "Invalid content type", http.StatusBadRequest)
		return
	}
	c.Next()
}

func AddAuthorizationContext(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Next()
		return
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Debug().Str("header", authHeader).Msg("Authorization header sent without Bearer prefix")
		c.Next()
		return
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := token_model.GetToken(accessToken)
	if err != nil {
		log.Debug().Err(err).Str("token", accessToken).Msg("Error getting token")
		c.Next()
		return
	}
	app, err := app_model.GetApp(token.ClientID)
	if err != nil {
		c.Next()
		return
	}

	c.Set(apicontext.AppContextKey, app)
	if token.UserID != nil {
		account, ok := account_model.GetAccountByID(*token.UserID)
		if !ok {
			c.Next()
			return
		}
		c.Set(apicontext.UserContextKey, account)
	}
	c.Next()
}

func RequireAppAuthorization(c *gin.Context) {
	app, ok := c.Get(apicontext.AppContextKey)
	if !ok || app == nil {
		http.Error(c.Writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	c.Next()
}

func RequireUserAuthentication(c *gin.Context) {
	account, ok := c.Get(apicontext.UserContextKey)
	if !ok || account == nil {
		http.Error(c.Writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	c.Next()
}

func Route(r *gin.RouterGroup) {
	r.Use(AddAuthorizationContext)

	v1 := r.Group("/v1")
	{
		v1.GET("/instance", instance_v1.Handler)

		v1.POST("/apps", apps.CreateAppHandler)
		v1.GET("/custom_emojis", instance.CustomEmojiHandler)
		v1.GET("/trends/tags", instance.TrendingTagsHandler)
		v1.GET("/trends/statuses", instance.TrendingStatusesHandler)
		v1.GET("/notifications", notifications.Handler)

		v1.GET("/timelines/public", timelines.PublicHandler)
		v1.GET("/timelines/home", timelines.HomeHandler)

		v1.GET("/accounts/lookup", accounts.LookupAccountHandler)
		v1.GET("/accounts/:id", accounts.GetAccountHandler)
		v1.GET("/accounts/:id/featured_tags", accounts.GetAccountHandler)
		v1.GET("/accounts/:id/statuses", accounts.StatusesHandler)

		appAuth := v1.Group("")
		{
			appAuth.Use(RequirePostJSONBody)
			appAuth.Use(RequireAppAuthorization)
			appAuth.POST("/accounts", accounts.RegisterAccountHandler)
			appAuth.GET("/apps/verify_credentials", apps.VerifyCredentialsHandler)
		}

		userAuth := v1.Group("")
		{
			userAuth.Use(RequirePostJSONBody)
			userAuth.Use(RequireUserAuthentication)
			userAuth.GET("/accounts/verify_credentials", accounts.VerifyCredentialsHandler)
		}
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/instance", instance_v2.Handler)
	}
}
