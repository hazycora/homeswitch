package app

import (
	"errors"

	"git.gay/h/homeswitch/config"
	"git.gay/h/homeswitch/crypto"
	"git.gay/h/homeswitch/db"

	"github.com/rs/zerolog/log"
)

var (
	ErrAppNotFound = errors.New("app not found")
	SystemApp      *App
)

func init() {
	db.Engine.Sync(new(App))
	SystemApp = &App{
		ID:      "0",
		Name:    "homeswitch",
		Website: config.ServerURL,
	}
	ok, err := db.Engine.Get(SystemApp)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting system app")
	}
	if !ok {
		err = CreateApp(SystemApp)
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating system app")
		}
	}
}

type App struct {
	ID           string   `json:"id" xorm:"'id' pk notnull unique"`
	Name         string   `json:"name" xorm:"'name' varchar(255) notnull"`
	RedirectURI  string   `json:"redirect_uris" xorm:"'redirect_uri'"`
	ClientID     string   `json:"client_id" xorm:"'client_id' varchar(255) notnull unique"`
	ClientSecret string   `json:"client_secret" xorm:"'client_secret' varchar(255) notnull"`
	Scopes       []string `json:"scopes" xorm:"'scopes'"`
	Website      string   `json:"website" xorm:"'website' varchar(255) notnull"`
	// TODO: Add VapidKey
}

func CreateApp(a *App) (err error) {
	if a.ID == "" {
		var id string
		id, err = db.RandomId()
		if err != nil {
			log.Error().Err(err).Msg("Error generating random ID")
			return
		}
		a.ID = id
	}
	clientId, err := crypto.RandomString(32)
	if err != nil {
		log.Error().Err(err).Msg("Error generating client secret")
		return
	}
	a.ClientID = clientId
	clientSecret, err := crypto.RandomString(32)
	if err != nil {
		log.Error().Err(err).Msg("Error generating client secret")
		return
	}
	a.ClientSecret = clientSecret
	_, err = db.Engine.Insert(a)
	return
}

func GetApp(clientId string) (app *App, err error) {
	app = &App{
		ClientID: clientId,
	}
	ok, err := db.Engine.Get(app)
	if err != nil {
		log.Error().Err(err).Msg("Error getting app")
		return
	}
	if !ok {
		err = ErrAppNotFound
		return
	}
	return
}
