package activitypub

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.gay/h/homeswitch/db"
	"git.gay/h/homeswitch/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

var (
	ErrActorNotFound = errors.New("actor not found")
)

func getActor(r *http.Request) (actor *models.Actor, err error) {
	username := chi.URLParam(r, "username")
	actor = &models.Actor{
		Username: username,
	}
	found, err := db.Engine.Get(actor)
	if err != nil {
		log.Error().Err(err).Msg("Error getting actor")
		return
	}
	if !found {
		err = ErrActorNotFound
		return
	}
	return
}

func ActorHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/activity+json")
	response := actor.ActivityPub()
	var body []byte
	body, err = json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling actor response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
