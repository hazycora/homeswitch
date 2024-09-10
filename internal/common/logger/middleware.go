package logger

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			wrapWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				if r := recover(); r != nil && r != http.ErrAbortHandler {
					logger.Error().Interface("recover", r).Bytes("stack", debug.Stack()).Msg("incoming_request_panic")
					wrapWriter.WriteHeader(http.StatusInternalServerError)
				}
				logger.Debug().Fields(map[string]interface{}{
					"remoteAddr": r.RemoteAddr,
					"path":       r.URL.Path,
					"proto":      r.Proto,
					"method":     r.Method,
					"userAgent":  r.UserAgent(),
					"statusCode": wrapWriter.Status(),
					"bytesIn":    r.ContentLength,
					"bytesOut":   wrapWriter.BytesWritten(),
				}).Msg("request")
			}()
			next.ServeHTTP(wrapWriter, r)
		}
		return http.HandlerFunc(fn)
	}
}
