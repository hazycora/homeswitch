package mastoapi

import (
	"net/http"

	// "git.gay/h/homeswitch/mastoapi/accounts"
	"git.gay/h/homeswitch/mastoapi/instance"
	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Get("/instance", instance.Handler)
		// r.Post("/accounts", accounts.RegisterAccountHandler)
	})
	return r
}
