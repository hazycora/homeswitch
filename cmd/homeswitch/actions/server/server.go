package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.gay/h/homeswitch/internal/router"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Start(ctx *cli.Context) error {
	host := ctx.String("host")
	port := ctx.String("port")
	addr := fmt.Sprintf("%s:%s", host, port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go router.Listen(addr)
	go func() {
		<-sigs
		done <- true
	}()

	<-done
	log.Info().Msg("Shutting down")

	return nil
}
