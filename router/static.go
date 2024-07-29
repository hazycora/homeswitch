package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetStaticRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Mount("/", http.FileServer(http.Dir("public")))
	return r
}
