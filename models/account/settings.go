package account

import (
	"encoding/json"

	"git.gay/h/homeswitch/models/status"
)

type AccountSettings struct {
	DefaultPrivacy   status.Privacy `json:"default_privacy"`
	DefaultSensitive bool           `json:"default_sensitive"`
	DefaultLanguage  *string        `json:"default_language"`
}

func (as *AccountSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(as)
}
