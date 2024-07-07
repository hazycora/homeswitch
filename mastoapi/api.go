package mastoapi

import (
	"context"
	"net/http"
	"strings"

	"git.gay/h/homeswitch/mastoapi/accounts"
	"git.gay/h/homeswitch/mastoapi/apicontext"
	"git.gay/h/homeswitch/mastoapi/apps"
	"git.gay/h/homeswitch/mastoapi/instance"
	app_model "git.gay/h/homeswitch/models/app"

	"github.com/go-chi/chi/v5"
)

func RequirePostJSONBody(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid content type", http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func RequireAppAuthorization(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}
		if !strings.HasSuffix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusBadRequest)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		app, err := app_model.GetAppByToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), apicontext.AppContextKey, app))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/apps", apps.CreateAppHandler)
		r.Get("/instance", instance.Handler)

		r.Group(func(r chi.Router) {
			r.Use(RequirePostJSONBody)
			r.Use(RequireAppAuthorization)
			r.Post("/apps/verify_credentials", apps.VerifyCredentialsHandler)
			r.Post("/accounts", accounts.RegisterAccountHandler)
		})
	})
	return r
}
