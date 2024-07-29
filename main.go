package main

import (
	"net/http"
	"os"

	actor_model "git.gay/h/homeswitch/models/actor"
	"git.gay/h/homeswitch/router"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	app := &cli.App{
		Name:  "homeswitch",
		Usage: "Start the homeswitch server",
		Action: func(*cli.Context) error {
			r := router.GetRouter()
			err := http.ListenAndServe("127.0.0.1:7983", r)
			return err
		},
		Commands: []*cli.Command{
			{
				Name: "add-account",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "username",
						Usage: "username for the account",
					},
					&cli.StringFlag{
						Name:  "password",
						Usage: "password for the account",
					},
					&cli.StringFlag{
						Name:  "email",
						Usage: "email for the account",
					},
				},
				Action: func(ctx *cli.Context) (err error) {
					username := ctx.String("username")
					email := ctx.String("email")
					password := ctx.String("password")
					if username == "" {
						log.Error().Msg("username unspecified")
					}
					if email == "" {
						log.Error().Msg("email unspecified")
					}
					if password == "" {
						log.Error().Msg("password unspecified")
					}
					actor := &actor_model.Actor{
						Username: username,
						Name:     &username,
						Email:    email,
					}
					err = actor_model.CreateActor(actor, ctx.String("password"))
					return
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Send()
	}
}
