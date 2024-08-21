package router

import (
	"net/http"
	"time"

	"git.gay/h/homeswitch/activitypub"
	actor_model "git.gay/h/homeswitch/models/actor"
	app_model "git.gay/h/homeswitch/models/app"
	token_model "git.gay/h/homeswitch/models/token"
	"git.gay/h/homeswitch/router/mastoapi"
	"git.gay/h/homeswitch/router/mastoapi/oauth"
	"git.gay/h/homeswitch/router/middleware"
	"git.gay/h/homeswitch/router/nodeinfo"
	"git.gay/h/homeswitch/router/tmpl"
	webfingerHandler "git.gay/h/homeswitch/webfinger/handler"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)
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
		user, ok := actor_model.ActorLogin(email, password)
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
	r.Mount("/", GetStaticRouter())
	return r
}
