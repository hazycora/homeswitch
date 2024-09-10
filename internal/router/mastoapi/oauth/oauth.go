package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.gay/h/homeswitch/internal/crypto"
	account_model "git.gay/h/homeswitch/internal/models/account"
	app_model "git.gay/h/homeswitch/internal/models/app"
	token_model "git.gay/h/homeswitch/internal/models/token"
	"git.gay/h/homeswitch/internal/router/mastoapi/form"
	"git.gay/h/homeswitch/internal/utils"

	"github.com/hazycora/go-mcache"
	"github.com/rs/zerolog/log"
)

var codeCache = mcache.New()

type AuthorizeForm struct {
	ResponseType string `json:"response_type" validate:"required,eq=code"`
	ClientID     string `json:"client_id" validate:"required"`
	RedirectURI  string `json:"redirect_uri" validate:"required"`
	Scope        string `json:"scope"`
	Lang         string `json:"lang"`
	// TODO: add support for force_login
}

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	requestForm := AuthorizeForm{
		ResponseType: r.Form.Get("response_type"),
		ClientID:     r.Form.Get("client_id"),
		RedirectURI:  r.Form.Get("redirect_uri"),
		Scope:        r.Form.Get("scope"),
		Lang:         r.Form.Get("lang"),
	}
	err := form.ValidateForm(requestForm)
	if err != nil {
		formError, ok := err.(form.FormError)
		if !ok {
			http.Error(w, "Error validating form", http.StatusInternalServerError)
			return
		}
		body, err := json.Marshal(formError)
		if err != nil {
			http.Error(w, "Error marshalling form error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	code, err := crypto.RandomString(32)
	if err != nil {
		http.Error(w, "Error generating code", http.StatusInternalServerError)
		return
	}

	accessToken, err := r.Cookie("access_token")
	if accessToken == nil || err != nil {
		// TODO: make a sign_in form here so that you don't need to already be signed in
		http.Error(w, "Unauthorized (Log in first at /auth/sign_in, then try again. Will make this easier later)", http.StatusUnauthorized)
		return
	}
	token, err := token_model.GetToken(accessToken.Value)
	if err != nil {
		http.Error(w, "Error getting token", http.StatusUnauthorized)
		return
	}
	codeCache.Set(code, *token.UserID, time.Hour)

	if requestForm.RedirectURI == "urn:ietf:wg:oauth:2.0:oob" {
		fmt.Fprintf(w, "Code: %s", code)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s?code=%s", requestForm.RedirectURI, url.QueryEscape(code)), http.StatusTemporaryRedirect)
}

type TokenForm struct {
	GrantType    string  `json:"grant_type" validate:"required"`
	Code         *string `json:"code"`
	ClientID     string  `json:"client_id" validate:"required"`
	ClientSecret string  `json:"client_secret" validate:"required"`
	RedirectURI  string  `json:"redirect_uri" validate:"required"`
	Scope        string  `json:"scope"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	utils.ParseForm(r)
	var formCode *string
	if code := r.FormValue("code"); code != "" {
		formCode = &code
	}
	requestForm := &TokenForm{
		GrantType:    r.FormValue("grant_type"),
		Code:         formCode,
		ClientID:     r.FormValue("client_id"),
		ClientSecret: r.FormValue("client_secret"),
		RedirectURI:  r.FormValue("redirect_uri"),
		Scope:        r.FormValue("scope"),
	}
	err := form.ValidateForm(*requestForm)
	if err != nil {
		formError, ok := err.(form.FormError)
		if !ok {
			http.Error(w, "Error validating form", http.StatusInternalServerError)
			return
		}
		body, err := json.Marshal(formError)
		if err != nil {
			http.Error(w, "Error marshalling form error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		log.Debug().Err(err).RawJSON("body", body).Msg("Invalid form at /oauth/token")
		return
	}

	var userId *string
	if requestForm.Code != nil {
		code := *requestForm.Code
		cachedUserId, ok := codeCache.Get(code)
		if !ok {
			log.Debug().Msg("Invalid code")
			http.Error(w, "Invalid code", http.StatusBadRequest)
			return
		}
		userIdString := cachedUserId.(string)
		userId = &userIdString
		codeCache.Remove(code)
	}

	app, err := app_model.GetApp(requestForm.ClientID)
	if err != nil || app.ClientSecret != requestForm.ClientSecret {
		log.Error().Err(err).Str("client_id", requestForm.ClientID).Msg("Failed to get app")
		if err == nil {
			log.Debug().Str("secret", app.ClientSecret).Str("claimed_secret", requestForm.ClientSecret).Msg("Secrets")
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := &token_model.Token{
		ClientID: r.Form.Get("client_id"),
	}
	if userId != nil {
		_, ok := account_model.GetAccountByID(*userId)
		if !ok {
			log.Debug().Str("user_id", *userId).Msg("Account not found")
			http.Error(w, "Account not found", http.StatusBadRequest)
		}
		token.UserID = userId
	}
	err = token_model.CreateToken(token)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"access_token": token.AccessToken,
		"token_type":   "Bearer",
		"scope":        strings.Join(token.Scopes, " "),
		"created_at":   token.CreatedAt,
	}
	body, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
