package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.gay/h/homeswitch/internal/router"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func listen(addr string) (err error) {
	r := router.GetRouter()
	log.Info().Msgf("Listening on http://%s", addr)
	err = http.ListenAndServe(addr, r)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("ListenAndServe error")
		return
	}
	return
}

func Start(ctx *cli.Context) error {
	host := ctx.String("host")
	port := ctx.String("port")
	addr := fmt.Sprintf("%s:%s", host, port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go listen(addr)
	go func() {
		<-sigs
		done <- true
	}()

	<-done
	log.Info().Msg("Shutting down")

	return nil
}
