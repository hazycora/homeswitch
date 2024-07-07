package apps

import (
	"encoding/json"
	"net/http"
	"strings"

	app_model "git.gay/h/homeswitch/models/app"
	"git.gay/h/homeswitch/router/mastoapi/apicontext"
	"git.gay/h/homeswitch/router/mastoapi/form"
	"github.com/rs/zerolog/log"
)

type CreateAppForm struct {
	ClientName  string   `json:"client_name" validate:"required"`
	RedirectURI string   `json:"redirect_uris" validate:"required"`
	Scopes      []string `json:"scopes"`
	Website     string   `json:"website"`
}

func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	requestForm := CreateAppForm{
		ClientName:  r.Form.Get("client_name"),
		RedirectURI: r.Form.Get("redirect_uris"),
		Scopes:      strings.Split(r.Form.Get("scopes"), " "),
		Website:     r.Form.Get("website"),
	}
	err := form.ValidateForm(requestForm)
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

	app := &app_model.App{
		Name:        requestForm.ClientName,
		RedirectURI: requestForm.RedirectURI,
		Scopes:      requestForm.Scopes,
		Website:     requestForm.Website,
	}

	err = app_model.CreateApp(app)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error creating app")
		http.Error(w, "Error creating app", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":            app.ID,
		"name":          app.Name,
		"website":       app.Website,
		"redirect_uri":  app.RedirectURI,
		"client_id":     app.ClientID,
		"client_secret": app.ClientSecret,
	}

	body, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error marshalling response")
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func VerifyCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(apicontext.AppContextKey).(*app_model.App)
	response := map[string]interface{}{
		"name":    app.Name,
		"website": app.Website,
	}
	body, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Str("path", r.URL.Path).Msg("Error marshalling response")
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
