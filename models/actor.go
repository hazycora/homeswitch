package models

import (
	"fmt"

	"git.gay/h/homeswitch/config"
	"git.gay/h/homeswitch/webfinger"
)

type Actor struct {
	ID         uint64  `xorm:"'id' pk autoincr"`
	Username   string  `json:"username" xorm:"varchar(25) notnull unique"`
	Name       *string `json:"name" xorm:"varchar(255) null"`
	Email      string  `json:"email"`
	Bio        *string `json:"note" xorm:"varchar(8096) null"`
	Created    int64   `json:"created" xorm:"'created'"`
	PrivateKey []byte  `json:"private_key" xorm:"notnull"`
	PublicKey  []byte  `json:"public_key" xorm:"notnull"`
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
