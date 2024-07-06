package main

import (
	"net/http"
	"os"

	"git.gay/h/homeswitch/activitypub"
	"git.gay/h/homeswitch/mastoapi"
	webfingerHandler "git.gay/h/homeswitch/webfinger/handler"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Server", "homeswitch (https://git.gay/h/homeswitch)")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("homeswitch!"))
	})
	r.Get("/.well-known/webfinger", webfingerHandler.Handler)
	r.Get("/@{username}", activitypub.ActorHandler)
	r.Post("/@{username}/inbox", activitypub.InboxHandler)
	r.Mount("/api", mastoapi.Router())
	http.ListenAndServe(":7983", r)
}
