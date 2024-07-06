package mastoapi

import (
	"net/http"

	// "git.gay/h/homeswitch/mastoapi/accounts"
	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		// r.Post("/accounts", accounts.RegisterAccountHandler)
	})
	return r
}
