package actor

import (
	"fmt"

	"git.gay/h/homeswitch/config"
	"git.gay/h/homeswitch/models/status"
	"git.gay/h/homeswitch/utils/marshaltime"
)

type Account struct {
	Actor
	URL             string            `json:"url"`
	URI             string            `json:"uri"`
	Group           bool              `json:"group"`
	AvatarURL       string            `json:"avatar"`
	AvatarStaticURL string            `json:"avatar_static"`
	HeaderURL       string            `json:"header"`
	HeaderStaticURL string            `json:"header_static"`
	LastStatusAt    *marshaltime.Date `json:"last_status_at"`
	Emoji           []interface{}     `json:"emojis"`
	Source          *AccountSource    `json:"source,omitempty"`
}

type AccountSource struct {
	Privacy             status.Privacy `json:"privacy"`
	Sensitive           bool           `json:"sensitive"`
	Language            *string        `json:"language"`
	Note                *string        `json:"note"`
	Fields              []Field        `json:"fields"`
	FollowRequestsCount int64          `json:"follow_requests_count"`
	HideCollections     *bool          `json:"hide_collections"`
	Discoverable        *bool          `json:"discoverable"`
	Indexable           bool           `json:"indexable"`
}

func (a *Actor) ToAccount(withSource bool) (account Account) {
	url := fmt.Sprintf("%s/@%s", config.ServerURL, a.Acct)
	avatarUrl := fmt.Sprintf("%s/static/missing_avatar.png", config.ServerURL)
	headerUrl := fmt.Sprintf("%s/static/missing_header.png", config.ServerURL)
	account = Account{
		Actor:           *a,
		URL:             url,
		URI:             url,
		Group:           false,
		Emoji:           []interface{}{},
		AvatarURL:       avatarUrl,
		AvatarStaticURL: avatarUrl,
		HeaderURL:       headerUrl,
		HeaderStaticURL: headerUrl,
	}
	if withSource {
		account.Source = &AccountSource{
			Privacy:             a.Settings.DefaultPrivacy,
			Sensitive:           a.Settings.DefaultSensitive,
			Language:            a.Settings.DefaultLanguage,
			Note:                a.Bio,
			Fields:              a.Fields,
			FollowRequestsCount: 0,
			HideCollections:     &a.HideCollections,
			Discoverable:        &a.Discoverable,
			Indexable:           a.Indexable,
		}
	}
	return
}
