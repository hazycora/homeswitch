package main

import (
	"net/http"
	"os"

	"git.gay/h/homeswitch/router"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	r := router.GetRouter()
	http.ListenAndServe(":7983", r)
}
