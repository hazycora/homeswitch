package middleware

import "net/http"

func AllowCORS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "accept, content-type, x-requested-with, authorization")
		w.Header().Add("Access-Control-Allow-Methods", "POST, PUT, DELETE, GET, PATCH, OPTIONS")
		w.Header().Add("Access-Control-Expose-Headers", "Link, X-RateLimit-Reset, X-RateLimit-Limit, X-RateLimit-Remaining, X-Request-Id")

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
