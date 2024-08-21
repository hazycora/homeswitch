package notifications

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("[]"))
}
