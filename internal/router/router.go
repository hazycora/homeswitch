package router

import (
	"errors"
	"net/http"

	"git.gay/h/homeswitch/internal/common/logger"
	"git.gay/h/homeswitch/internal/router/activitypub"
	"git.gay/h/homeswitch/internal/router/mastoapi"
	"git.gay/h/homeswitch/internal/router/mastoapi/oauth"
	"git.gay/h/homeswitch/internal/router/middleware"
	"git.gay/h/homeswitch/internal/router/nodeinfo"
	webfingerHandler "git.gay/h/homeswitch/internal/webfinger/handler"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = logger.Logger
}

func GetRouter() http.Handler {
	r := gin.New()
	r.Use(logger.Middleware)
	r.Use(func(c *gin.Context) {
		c.Header("Server", "homeswitch (https://git.gay/h/homeswitch)")
		c.Next()
	})
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.Writer.WriteString(`<style>:root {color-scheme: light dark; font-family: system-ui, sans-serif;}</style>
		<a href="https://homeswit.ch" target="_blank">homeswitch</a> - a work-in-progress fediverse server.`)
	})
	r.GET("/.well-known/webfinger", webfingerHandler.Handler)
	r.GET("/.well-known/nodeinfo", nodeinfo.WellKnownHandler)
	r.GET("/nodeinfo/2.0", nodeinfo.Handler)
	r.GET("/nodeinfo/2.0.json", nodeinfo.Handler)

	r.GET("/@{username}", activitypub.ActorHandler)
	r.POST("/@{username}/inbox", activitypub.InboxHandler)

	apiRoute := r.Group("/api")
	mastoapi.Route(apiRoute)

	corsGroup := r.Group("")
	{
		corsGroup.Use(middleware.AllowCORS)

		corsGroup.GET("/oauth/authorize", oauth.AuthorizeHandler)
		corsGroup.POST("/oauth/token", oauth.TokenHandler)
	}

	r.POST("/auth/sign_in")
	r.Static("/system/static", "static")
	return r
}

func Listen(addr string) (err error) {
	router := GetRouter()
	log.Info().Msgf("Listening on http://%s", addr)
	err = http.ListenAndServe(addr, router)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("ListenAndServe error")
		return
	}
	return
}
