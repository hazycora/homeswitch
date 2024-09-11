package router

import (
	"errors"
	"net/http"
	"time"

	"git.gay/h/homeswitch/internal/activitypub"
	"git.gay/h/homeswitch/internal/common/logger"
	account_model "git.gay/h/homeswitch/internal/models/account"
	app_model "git.gay/h/homeswitch/internal/models/app"
	token_model "git.gay/h/homeswitch/internal/models/token"
	"git.gay/h/homeswitch/internal/router/mastoapi"
	"git.gay/h/homeswitch/internal/router/mastoapi/oauth"
	"git.gay/h/homeswitch/internal/router/middleware"
	"git.gay/h/homeswitch/internal/router/nodeinfo"
	"git.gay/h/homeswitch/internal/router/tmpl"
	webfingerHandler "git.gay/h/homeswitch/internal/webfinger/handler"

	"github.com/chi-middleware/proxy"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = logger.Logger
}

func GetRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(proxy.ForwardedHeaders())
	r.Use(logger.Middleware)
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Server", "homeswitch (https://git.gay/h/homeswitch)")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`<style>:root {color-scheme: light dark; font-family: system-ui, sans-serif;}</style>
		<a href="https://homeswit.ch" target="_blank">homeswitch</a> - a work-in-progress fediverse server.`))
	})
	r.Get("/.well-known/webfinger", webfingerHandler.Handler)
	r.Get("/.well-known/nodeinfo", nodeinfo.WellKnownHandler)
	r.Get("/nodeinfo/2.0", nodeinfo.Handler)
	r.Get("/nodeinfo/2.0.json", nodeinfo.Handler)

	r.Get("/@{username}", activitypub.ActorHandler)
	r.Post("/@{username}/inbox", activitypub.InboxHandler)

	r.Mount("/api", mastoapi.Router())

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowCORS)

		r.Get("/oauth/authorize", oauth.AuthorizeHandler)
		r.Post("/oauth/token", oauth.TokenHandler)
		r.Get("/auth/sign_in", func(w http.ResponseWriter, r *http.Request) {
			tmpl.Render(w, "auth/signin", nil)
		})
	})
	r.Post("/auth/sign_in", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := r.Form.Get("user[email]")
		password := r.Form.Get("user[password]")
		user, ok := account_model.AccountLogin(email, password)
		if !ok {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		token := &token_model.Token{
			ClientID:  app_model.SystemApp.ClientID,
			TokenType: "Bearer",
			UserID:    &user.ID,
		}
		err := token_model.CreateToken(token)
		if err != nil {
			http.Error(w, "Error creating token", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token.AccessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteDefaultMode,
			MaxAge:   int(time.Hour.Seconds() * 24 * 30),
		})
		w.Write([]byte("Signed in!"))
	})
	r.Mount("/system/static", http.StripPrefix("/system/static/", http.FileServer(http.Dir("static"))))
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
