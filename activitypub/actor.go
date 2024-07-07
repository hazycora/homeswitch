package activitypub

import (
	"encoding/json"
	"net/http"

	actor_model "git.gay/h/homeswitch/models/actor"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	actor, ok := actor_model.GetActorByUsername(username)
	if !ok {
		http.Error(w, "Actor not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/activity+json")
	response := actor.ActivityPub()
	var body []byte
	body, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling actor response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
