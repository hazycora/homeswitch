package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"git.gay/h/homeswitch/cmd/homeswitch/actions/account"
	"git.gay/h/homeswitch/cmd/homeswitch/actions/server"
	"git.gay/h/homeswitch/internal/common/logger"
)

func init() {
	log.Logger = logger.Logger
}

func main() {
	app := &cli.App{
		Name:  "homeswitch",
		Usage: "Start the homeswitch server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Usage: "Host for the server",
				Value: "127.0.0.1",
			},
			&cli.StringFlag{
				Name:  "port",
				Usage: "Port for the server",
				Value: "7983",
			},
		},
		Action: server.Start,
		Commands: []*cli.Command{
			account.Add,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Send()
	}
}
