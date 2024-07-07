package accounts

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"git.gay/h/homeswitch/crypto"
	"git.gay/h/homeswitch/db"
	"git.gay/h/homeswitch/mastoapi/form"
	"git.gay/h/homeswitch/models/actor"
	"github.com/rs/zerolog/log"

	"github.com/alexedwards/argon2id"
)

type RegisterAccountForm struct {
	Username  string `json:"username" validate:"required,username,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Agreement bool   `json:"agreement" validate:"required"`
	Locale    string `json:"locale" validate:"required"`
	Reason    string `json:"reason"`
}

func RegisterAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error reading body")
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	var requestForm RegisterAccountForm
	err = json.Unmarshal(body, &requestForm)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error unmarshalling JSON")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err = form.ValidateForm(requestForm)
	if err != nil {
		formError, ok := err.(form.FormError)
		if !ok {
			log.Error().Err(err).Str("path", r.URL.Path).Msg("Error validating form")
			http.Error(w, "Error validating form", http.StatusInternalServerError)
			return
		}
		log.Debug().Err(formError).Str("path", r.URL.Path).Msg("Received invalid form")
		body, err := json.Marshal(formError)
		if err != nil {
			log.Error().Err(err).Str("path", r.URL.Path).Msg("Error marshalling form error")
			http.Error(w, "Error marshalling form error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(body)
		return
	}

	id, err := db.RandomId()
	if err != nil {
		log.Error().Err(err).Str("username", requestForm.Username).Msg("Error generating random ID")
		http.Error(w, "Error generating random ID", http.StatusInternalServerError)
		return
	}
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Error().Err(err).Str("username", requestForm.Username).Msg("Error generating key pair")
		http.Error(w, "Error generating key pair", http.StatusInternalServerError)
		return
	}

	hash, err := argon2id.CreateHash(requestForm.Password, argon2id.DefaultParams)
	if err != nil {
		log.Error().Err(err).Str("username", requestForm.Username).Msg("Error hashing password")
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	actor := actor.Actor{
		ID:           id,
		Username:     requestForm.Username,
		Name:         &requestForm.Username,
		Email:        requestForm.Email,
		Created:      time.Now().Unix(),
		PrivateKey:   string(privateKey),
		PublicKey:    string(publicKey),
		PasswordHash: hash,
	}
	_, err = db.Engine.Insert(&actor)
	if err != nil {
		log.Error().Err(err).Str("username", requestForm.Username).Msg("Error inserting actor")
		http.Error(w, "Error inserting actor", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
