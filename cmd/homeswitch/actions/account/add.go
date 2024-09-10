package account

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	account_model "git.gay/h/homeswitch/internal/models/account"
)

var Add = &cli.Command{
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
		account := &account_model.Account{
			Username: username,
			Name:     &username,
			Email:    email,
		}
		err = account_model.CreateAccount(account, ctx.String("password"))
		return
	},
}
