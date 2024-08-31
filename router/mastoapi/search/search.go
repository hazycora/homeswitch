package search

import (
	"net/http"
	"strings"

	account_model "git.gay/h/homeswitch/models/account"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "param is missing or the value is empty", http.StatusBadRequest)
	}

	query = strings.TrimSpace(query)

	if account_model.AcctRegExp.MatchString(query) {
		account_model.GetAccountByUsername(query)
	}
}
