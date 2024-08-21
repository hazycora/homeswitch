package accounts

import (
	"encoding/json"
	"io"
	"net/http"

	"git.gay/h/homeswitch/config"
	actor_model "git.gay/h/homeswitch/models/actor"
	"git.gay/h/homeswitch/router/mastoapi/apicontext"
	"git.gay/h/homeswitch/router/mastoapi/form"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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
	if !config.RegistrationsEnabled {
		http.Error(w, "Registrations not enabled.", http.StatusForbidden)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
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

	actor := &actor_model.Actor{
		Username: requestForm.Username,
		Name:     &requestForm.Username,
		Email:    requestForm.Email,
	}
	err = actor_model.CreateActor(actor, requestForm.Password)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func VerifyCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	actor := r.Context().Value(apicontext.UserContextKey).(*actor_model.Actor)
	body, err := json.Marshal(actor.ToAccount(true))
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error marshalling response")
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(body)
	log.Debug().Any("body", string(body)).Msg("verify_credentials, sob")
}

func GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "id")
	actor, ok := actor_model.GetActorByID(accountId)
	if !ok {
		// TODO: error should be same as Mastodon
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}
	account := actor.ToAccount(false)

	body, err := json.Marshal(account)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}

func LookupAccountHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	acct := r.Form.Get("acct")
	actor, ok := actor_model.GetActorByUsername(acct)
	if !ok {
		// TODO: error should be same as Mastodon
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}
	account := actor.ToAccount(false)

	body, err := json.Marshal(account)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}

// TODO: implement
func FeaturedTagsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("[]"))
}

// TODO: implement
func StatusesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("[]"))
}
