package actor

import (
	"errors"

	"git.gay/h/homeswitch/crypto"
	"git.gay/h/homeswitch/db"
	"git.gay/h/homeswitch/utils/marshaltime"
	"github.com/alexedwards/argon2id"
	"github.com/rs/zerolog/log"
)

func CreateActor(actor *Actor, password string) (err error) {
	id, err := db.RandomId()
	if err != nil {
		errors.Join(err, errors.New("generating random ID failed"))
		return
	}
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		errors.Join(err, errors.New("generating key pair failed"))
		return
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		errors.Join(err, errors.New("hashing password failed"))
		return
	}

	actor.ID = id
	actor.Acct = actor.Username
	actor.Created = marshaltime.Now()
	actor.PrivateKey = string(privateKey)
	actor.PublicKey = string(publicKey)
	actor.PasswordHash = hash
	actor.Fields = []Field{}

	_, err = db.Engine.Insert(actor)
	if err != nil {
		log.Error().Err(err).Str("username", actor.Username).Msg("Error inserting actor")
		return
	}
	return
}
