package status

import "encoding/json"

type Privacy int

const (
	PrivacyPublic Privacy = iota
	PrivacyUnlisted
	PrivacyFollowers
	PrivacyMentioned
)

func (p Privacy) String() string {
	switch p {
	case PrivacyPublic:
		return "public"
	case PrivacyUnlisted:
		return "unlisted"
	case PrivacyFollowers:
		return "private"
	case PrivacyMentioned:
		return "direct"
	}
	panic("unknown privacy type")
}
func (p Privacy) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
