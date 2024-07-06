package activitypub

import (
	"errors"
	"io"
	"net/http"

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
	actor, err := getActor(r)
	if err != nil {
		if errors.Is(err, ErrActorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("Error getting actor")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug().Str("body", string(body)).Str("for_actor", actor.Username).Msg("Inbox event, got body")
	// TODO: actually handle inbox events
}
