package utils

import (
	"encoding/json"
	"net/http"
)

// Handles application/json request bodies
func ParseForm(r *http.Request) (err error) {
	r.ParseForm()
	if r.Header.Get("Content-Type") == "application/json" {
		body := make(map[string]string)
		err = json.NewDecoder(r.Body).Decode(&body)
		for key, value := range body {
			r.Form[key] = []string{value}
		}
	}
	return
}