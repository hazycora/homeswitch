package instance

import (
	"encoding/json"
	"net/http"
)

func CustomEmojiHandler(w http.ResponseWriter, r *http.Request) {
	emoji := []interface{}{}
	body, err := json.Marshal(emoji)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}