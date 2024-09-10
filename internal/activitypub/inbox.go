package activitypub

import (
	"io"
	"net/http"

	account_model "git.gay/h/homeswitch/internal/models/account"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error reading body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := chi.URLParam(r, "username")
	account, ok := account_model.GetAccountByUsername(username)
	if !ok {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	log.Debug().Str("body", string(body)).Str("for_account", account.Username).Msg("Inbox event, got body")
	// TODO: actually handle inbox events
}
