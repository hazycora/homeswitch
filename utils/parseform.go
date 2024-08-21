package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alecthomas/units"
	"github.com/rs/zerolog/log"
)

func isMultipartForm(r *http.Request) bool {
	if r.Header.Get("Content-Type") == "multipart/form-data" {
		return true
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data;") {
		return true
	}
	return false
}

// Handles application/json request bodies
func ParseForm(r *http.Request) (err error) {
	if isMultipartForm(r) {
		log.Debug().Msg("Multipart form!")
		r.ParseMultipartForm(int64(units.Mebibyte))
	} else {
		r.ParseForm()
	}
	if r.Header.Get("Content-Type") == "application/json" {
		body := make(map[string]string)
		err = json.NewDecoder(r.Body).Decode(&body)
		for key, value := range body {
			r.Form[key] = []string{value}
		}
	}
	return
}
