package logger

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger
var Middleware func(next http.Handler) http.Handler

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	Middleware = LoggerMiddleware(Logger)
}
