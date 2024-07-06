package accounts

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"git.gay/h/homeswitch/crypto"
	"git.gay/h/homeswitch/db"
	"git.gay/h/homeswitch/models"
	"github.com/rs/zerolog/log"
)

type RegisterAccountForm struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Agreement bool   `json:"agreement"`
	Locale    string `json:"locale"`
	Reason    string `json:"reason"`
}

func RegisterAccountHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Require authentication
	// TODO: send proper error codes

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Error().Str("path", r.URL.Path).Str("content-type", contentType).Msg("Invalid content type")
		http.Error(w, "Invalid content type", http.StatusBadRequest)
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error reading body")
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	var form RegisterAccountForm
	err = json.Unmarshal(body, &form)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error unmarshalling JSON")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if !form.Agreement {
		log.Error().Str("username", form.Username).Msg("User did not agree to terms")
		http.Error(w, "User did not agree to terms", http.StatusUnauthorized)
		return
	}
	id, err := db.RandomId()
	if err != nil {
		log.Error().Err(err).Str("username", form.Username).Msg("Error generating random ID")
		http.Error(w, "Error generating random ID", http.StatusInternalServerError)
		return
	}
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Error().Err(err).Str("username", form.Username).Msg("Error generating key pair")
		http.Error(w, "Error generating key pair", http.StatusInternalServerError)
		return
	}
	actor := models.Actor{
		ID:         id,
		Username:   form.Username,
		Name:       &form.Username,
		Email:      form.Email,
		Created:    time.Now().Unix(),
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
	_, err = db.Engine.Insert(&actor)
	if err != nil {
		log.Error().Err(err).Str("username", form.Username).Msg("Error inserting actor")
		http.Error(w, "Error inserting actor", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
