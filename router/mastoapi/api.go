package mastoapi

import (
	"context"
	"net/http"
	"strings"

	account_model "git.gay/h/homeswitch/models/account"
	app_model "git.gay/h/homeswitch/models/app"
	token_model "git.gay/h/homeswitch/models/token"
	"git.gay/h/homeswitch/router/mastoapi/accounts"
	"git.gay/h/homeswitch/router/mastoapi/apicontext"
	"git.gay/h/homeswitch/router/mastoapi/apps"
	"git.gay/h/homeswitch/router/mastoapi/instance"
	instance_v1 "git.gay/h/homeswitch/router/mastoapi/instance/v1"
	instance_v2 "git.gay/h/homeswitch/router/mastoapi/instance/v2"
	"git.gay/h/homeswitch/router/mastoapi/notifications"
	"git.gay/h/homeswitch/router/mastoapi/timelines"
	"git.gay/h/homeswitch/router/middleware"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
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

func AddAuthorizationContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Debug().Str("header", authHeader).Msg("Authorization header sent without Bearer prefix")
			h.ServeHTTP(w, r)
			return
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := token_model.GetToken(accessToken)
		if err != nil {
			log.Debug().Err(err).Str("token", accessToken).Msg("Error getting token")
			h.ServeHTTP(w, r)
			return
		}
		app, err := app_model.GetApp(token.ClientID)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), apicontext.AppContextKey, app))
		if token.UserID != nil {
			account, ok := account_model.GetAccountByID(*token.UserID)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), apicontext.UserContextKey, account))
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func RequireAppAuthorization(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app, ok := r.Context().Value(apicontext.AppContextKey).(*app_model.App)
		if !ok || app == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func RequireUserAuthentication(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		account, ok := r.Context().Value(apicontext.UserContextKey).(*account_model.Account)
		if !ok || account == nil {
			log.Debug().Msg("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(chi_middleware.StripSlashes)
	r.Use(middleware.AllowCORS)

	r.Route("/v1", func(r chi.Router) {
		r.Use(AddAuthorizationContext)
		r.Post("/apps", apps.CreateAppHandler)
		r.Get("/custom_emojis", instance.CustomEmojiHandler)
		r.Get("/instance", instance_v1.Handler)
		r.Group(func(r chi.Router) {
			r.Get("/notifications", notifications.Handler)
		})

		r.Group(func(r chi.Router) {
			r.Get("/timelines/home", timelines.HomeHandler)
		})

		r.Group(func(r chi.Router) {
			r.Get("/accounts/lookup", accounts.LookupAccountHandler)
			r.Get("/accounts/{id}", accounts.GetAccountHandler)
			r.Get("/accounts/{id}/featured_tags", accounts.GetAccountHandler)
			r.Get("/accounts/{id}/statuses", accounts.StatusesHandler)
		})

		r.Group(func(r chi.Router) {
			r.Use(RequirePostJSONBody)
			r.Use(RequireAppAuthorization)
			r.Post("/accounts", accounts.RegisterAccountHandler)
			r.Get("/apps/verify_credentials", apps.VerifyCredentialsHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(RequirePostJSONBody)
			r.Use(RequireUserAuthentication)
			r.Get("/accounts/verify_credentials", accounts.VerifyCredentialsHandler)
		})
	})

	r.Route("/v2", func(r chi.Router) {
		r.Use(AddAuthorizationContext)
		r.Get("/instance", instance_v2.Handler)
	})
	return r
}
