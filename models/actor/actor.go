package actor

import (
	"errors"
	"fmt"

	"git.gay/h/homeswitch/config"
	"git.gay/h/homeswitch/db"
	"git.gay/h/homeswitch/utils/marshaltime"
	"git.gay/h/homeswitch/webfinger"
	"github.com/alexedwards/argon2id"
	"github.com/rs/zerolog/log"
)

var (
	ErrActorNotFound = errors.New("actor not found")
	EmptyRole        = map[string]interface{}{
		"id":          "-99",
		"name":        "",
		"permissions": "0",
		"color":       "",
		"highlighted": false,
	}
)

func init() {
	db.Engine.Sync(new(Actor))
}

type Actor struct {
	ID              string           `json:"id" xorm:"'id' pk notnull unique"`
	Username        string           `json:"username" xorm:"'username' varchar(25) notnull"`
	Acct            string           `json:"acct" xorm:"'acct' varchar(255) notnull unique"`
	Name            *string          `json:"display_name" xorm:"'name' varchar(255) null"`
	Email           string           `json:"-" xorm:"'email'"`
	AvatarID        *string          `json:"-" xorm:"'avatar' null"`
	HeaderID        *string          `json:"-" xorm:"'header' null"`
	Bio             *string          `json:"note" xorm:"'bio' varchar(8096) null"`
	Created         marshaltime.Time `json:"created_at" xorm:"'created' created"`
	PrivateKey      string           `json:"-" xorm:"'private_key' text notnull"`
	PublicKey       string           `json:"-" xorm:"'public_key' text notnull"`
	PasswordHash    string           `json:"-" xorm:"'password_hash' varchar(128) notnull"`
	Locked          bool             `json:"locked"`
	Bot             bool             `json:"bot"`
	Discoverable    bool             `json:"discoverable"`
	Indexable       bool             `json:"indexable"`
	NoIndex         bool             `json:"noindex"`
	HideCollections bool             `json:"hide_collections"`
	FollowersCount  int64            `json:"followers_count"`
	FollowingCount  int64            `json:"following_count"`
	StatusesCount   int64            `json:"statuses_count"`
	Fields          []Field          `json:"fields"`
	Settings        ActorSettings    `json:"-" xorm:"jsonb"`
}

type Field struct {
	Name       string            `json:"name"`
	Value      string            `json:"value"`
	VerifiedAt *marshaltime.Time `json:"verified_at"`
}

func (a *Actor) TableName() string {
	return "actor"
}

func (a *Actor) Webfinger() webfinger.Webfinger {
	return webfinger.Webfinger{
		Subject: fmt.Sprintf("acct:%s@%s", a.Username, config.ServerName),
		Aliases: []string{
			fmt.Sprintf("https://%s/@%s", config.ServerName, a.Username),
		},
		Links: []webfinger.WebfingerLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("https://%s/@%s", config.ServerName, a.Username),
			},
		},
	}
}

func (a *Actor) ActivityPub() map[string]interface{} {
	return map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		"type":              "Person",
		"id":                fmt.Sprintf("https://%s/@%s", config.ServerName, a.Username),
		"preferredUsername": a.Username,
		"name":              a.Name,
		"url":               fmt.Sprintf("https://%s/@%s", config.ServerName, a.Username),
		"summary":           a.Bio,
		"inbox":             fmt.Sprintf("https://%s/@%s/inbox", config.ServerName, a.Username),
		"publicKey": map[string]interface{}{
			"id":           fmt.Sprintf("https://%s/@%s#main-key", config.ServerName, a.Username),
			"owner":        fmt.Sprintf("https://%s/@%s", config.ServerName, a.Username),
			"publicKeyPem": string(a.PublicKey),
		},
	}
}

func GetActorByUsername(username string) (actor *Actor, ok bool) {
	actor = &Actor{
		Username: username,
	}
	exists, err := db.Engine.Get(actor)
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error getting actor by username")
		return
	}
	if !exists {
		return
	}
	ok = true
	return
}

func GetActorByID(id string) (actor *Actor, ok bool) {
	actor = &Actor{
		ID: id,
	}
	exists, err := db.Engine.Get(actor)
	if err != nil {
		log.Err(err).Str("id", id).Msg("Error getting actor by ID")
		return
	}
	if !exists {
		return
	}
	ok = true
	return
}

func ActorLogin(email string, password string) (actor *Actor, ok bool) {
	actor = &Actor{
		Email: email,
	}
	exists, err := db.Engine.Get(actor)
	if err != nil {
		log.Err(err).Str("email", email).Msg("Error getting actor by email")
		return
	}
	if !exists {
		return
	}
	match, err := argon2id.ComparePasswordAndHash(password, actor.PasswordHash)
	if err != nil {
		log.Err(err).Str("email", email).Str("username", actor.Username).Msg("Error comparing password and hash")
		return
	}
	if !match {
		return
	}
	ok = true
	return
}

func GetLocalActorCount() (count int64, err error) {
	count, err = db.Engine.Count(new(Actor))
	return
}
