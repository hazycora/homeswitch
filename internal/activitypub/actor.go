package activitypub

import (
	"encoding/json"
	"net/http"

	account_model "git.gay/h/homeswitch/internal/models/account"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	account, ok := account_model.GetAccountByUsername(username)
	if !ok {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/activity+json")
	response := account.ActivityPub()
	var body []byte
	body, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling account response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
