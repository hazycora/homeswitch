package actor

import (
	"encoding/json"

	"git.gay/h/homeswitch/models/status"
)

type ActorSettings struct {
	DefaultPrivacy   status.Privacy `json:"default_privacy"`
	DefaultSensitive bool           `json:"default_sensitive"`
	DefaultLanguage  *string        `json:"default_language"`
}

func (as *ActorSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(as)
}
