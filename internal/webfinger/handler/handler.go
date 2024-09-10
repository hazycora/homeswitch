package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"git.gay/h/homeswitch/internal/config"
	account_model "git.gay/h/homeswitch/internal/models/account"
	"github.com/rs/zerolog/log"
)

var (
	reResource = regexp.MustCompile(`^acct:([^@]+)@(.+)$`)
)

func Handler(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	submatches := reResource.FindStringSubmatch(resource)
	username := submatches[1]
	instance := submatches[2]
	if instance != config.ServerName {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Debug().Str("username", username).Any("instance", instance).Msg("Webfinger request")
	account, ok := account_model.GetAccountByUsername(username)
	if !ok {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/jrd+json")
	webfingerResponse, err := json.Marshal(account.Webfinger())
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling webfinger response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(webfingerResponse)
}
