package token

import (
	"time"

	"git.gay/h/homeswitch/crypto"
	"git.gay/h/homeswitch/db"
	"git.gay/h/homeswitch/models/app"
)

type Token struct {
	ClientID    string   `json:"client_id"`
	TokenType   string   `json:"token_type"`
	Scopes      []string `json:"scope"`
	UserID      *string  `json:"user_id" xorm:"null"`
	AccessToken string   `json:"access_token" xorm:"varchar(255) notnull unique"`
	CreatedAt   int64    `json:"created_at"`
}

func init() {
	db.Engine.Sync(new(Token))
}

func CreateToken(t *Token) (err error) {
	if t.Scopes == nil || len(t.Scopes) == 0 {
		t.Scopes = []string{"read"}
	}
	_, err = app.GetApp(t.ClientID)
	if err != nil {
		return
	}
	t.AccessToken, err = crypto.RandomString(32)
	if err != nil {
		return
	}
	t.CreatedAt = time.Now().Unix()
	t.TokenType = "Bearer"
	db.Engine.Insert(t)
	return
}

func GetToken(accessToken string) (t *Token, err error) {
	t = &Token{AccessToken: accessToken}
	ok, err := db.Engine.Get(t)
	if !ok {
		t = nil
	}
	return
}
