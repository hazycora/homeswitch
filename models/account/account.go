package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"git.gay/h/homeswitch/config"
	"git.gay/h/homeswitch/db"
	"git.gay/h/homeswitch/utils/marshaltime"
	"git.gay/h/homeswitch/webfinger"
	"github.com/alexedwards/argon2id"
	"github.com/rs/zerolog/log"
)

var (
	AcctRegExp         = regexp.MustCompile(`(?i)@?([a-z0-9\-\_]+)@([a-z0-9\-.]+)`)
	ErrAccountNotFound = errors.New("account not found")
	EmptyRole          = map[string]interface{}{
		"id":          "-99",
		"name":        "",
		"permissions": "0",
		"color":       "",
		"highlighted": false,
	}
)

func init() {
	db.Engine.Sync(new(Account))
}

type Account struct {
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
	Settings        AccountSettings  `json:"-" xorm:"jsonb"`
}

type Field struct {
	Name       string            `json:"name"`
	Value      string            `json:"value"`
	VerifiedAt *marshaltime.Time `json:"verified_at"`
}

func (a *Account) TableName() string {
	return "account"
}

func (a *Account) Webfinger() webfinger.Webfinger {
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

func (a *Account) ActivityPub() map[string]interface{} {
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

func lookupAccount(acct string) (account *Account, err error) {
	// TODO: support host-meta (eg: https://besties.house/.well-known/host-meta)
	submatches := AcctRegExp.FindStringSubmatch(acct)
	instance := submatches[2]
	if acct[0] == '@' {
		acct = acct[1:]
	}
	webfingerResponse, err := webfinger.LookupResource(instance, fmt.Sprintf("acct:%s", acct))
	if err != nil {
		return
	}
	link, ok := webfingerResponse.GetLink("application/activity+json")
	if !ok {
		err = fmt.Errorf("No application/activity+json link")
	}
	accountResponse, err := http.Get(link.Href)
	if err != nil {
		return
	}
	defer accountResponse.Body.Close()
	account = &Account{}
	err = json.NewDecoder(accountResponse.Body).Decode(account)
	if err != nil {
		return
	}
	err = fmt.Errorf("Not finished implementing")
	return
}

func GetAccountByUsername(username string) (account *Account, ok bool) {
	account = &Account{}
	isAcct := AcctRegExp.MatchString(username)
	if isAcct {
		var acct string
		if username[0] == '@' {
			acct = username[1:]
		} else {
			acct = username
		}
		account.Acct = acct
		lookupAccount(acct)
	} else {
		account.Username = username
	}
	exists, err := db.Engine.Get(account)
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error getting account by username")
		return
	}
	if !exists {
		return
	}
	ok = true
	return
}

func GetAccountByID(id string) (account *Account, ok bool) {
	account = &Account{
		ID: id,
	}
	exists, err := db.Engine.Get(account)
	if err != nil {
		log.Err(err).Str("id", id).Msg("Error getting account by ID")
		return
	}
	if !exists {
		return
	}
	ok = true
	return
}

func AccountLogin(email string, password string) (account *Account, ok bool) {
	account = &Account{
		Email: email,
	}
	exists, err := db.Engine.Get(account)
	if err != nil {
		log.Err(err).Str("email", email).Msg("Error getting account by email")
		return
	}
	if !exists {
		return
	}
	match, err := argon2id.ComparePasswordAndHash(password, account.PasswordHash)
	if err != nil {
		log.Err(err).Str("email", email).Str("username", account.Username).Msg("Error comparing password and hash")
		return
	}
	if !match {
		return
	}
	ok = true
	return
}

func GetLocalAccountCount() (count int64, err error) {
	count, err = db.Engine.Count(new(Account))
	return
}
