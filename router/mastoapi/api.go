package mastoapi

import (
	"context"
	"net/http"
	"strings"

	actor_model "git.gay/h/homeswitch/models/actor"
	app_model "git.gay/h/homeswitch/models/app"
	token_model "git.gay/h/homeswitch/models/token"
	"git.gay/h/homeswitch/router/mastoapi/accounts"
	"git.gay/h/homeswitch/router/mastoapi/apicontext"
	"git.gay/h/homeswitch/router/mastoapi/apps"
	"git.gay/h/homeswitch/router/mastoapi/instance"
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
		log.Debug().Str("token ClientID", token.ClientID).Msg("token.ClientID != nil")
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), apicontext.AppContextKey, app))
		if token.UserID != nil {
			log.Debug().Str("token UserID", *token.UserID).Msg("token.UserID != nil")
			actor, ok := actor_model.GetActorByID(*token.UserID)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), apicontext.UserContextKey, actor))
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
		actor, ok := r.Context().Value(apicontext.UserContextKey).(*actor_model.Actor)
		if !ok || actor == nil {
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
		r.Get("/instance", instance.Handler)

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
	return r
}
