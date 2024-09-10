package account

import (
	"errors"

	"git.gay/h/homeswitch/internal/crypto"
	"git.gay/h/homeswitch/internal/db"
	"git.gay/h/homeswitch/internal/utils/marshaltime"
	"github.com/alexedwards/argon2id"
	"github.com/rs/zerolog/log"
)

func CreateAccount(account *Account, password string) (err error) {
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

	account.ID = id
	account.Acct = account.Username
	account.Created = marshaltime.Now()
	account.PrivateKey = string(privateKey)
	account.PublicKey = string(publicKey)
	account.PasswordHash = hash
	account.Fields = []Field{}

	_, err = db.Engine.Insert(account)
	if err != nil {
		log.Error().Err(err).Str("username", account.Username).Msg("Error inserting account")
		return
	}
	return
}
