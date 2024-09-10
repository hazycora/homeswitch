package nodeinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.gay/h/homeswitch/internal/config"
	"git.gay/h/homeswitch/internal/version"
	"github.com/rs/zerolog/log"
)

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

var WellKnownLinks []byte

func init() {
	var links = []Link{
		{
			Rel:  "http://nodeinfo.diaspora.software/ns/schema/2.0",
			Href: fmt.Sprintf("%s/nodeinfo/2.0", config.ServerURL),
		},
	}
	var err error
	WellKnownLinks, err = json.Marshal(links)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to marshall .well-known/nodeinfo links")
	}
}

type NodeInfo struct {
	Version   string   `json:"version"`
	Software  Software `json:"software"`
	Protocols []string `json:"protocols"`
	Services  Services `json:"services"`
	Usage     struct {
		Users struct {
			Total          int64 `json:"total"`
			ActiveMonth    int64 `json:"activeMonth"`
			ActiveHalfYear int64 `json:"activeHalfYear"`
		} `json:"users"`
		LocalPosts int64 `json:"localPosts"`
	} `json:"usage"`
	OpenRegistration bool     `json:"openRegistration"`
	Metadata         Metadata `json:"metadata"`
}

type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Services struct {
	Inbound  []string `json:"outbound"`
	Outbound []string `json:"inbound"`
}

type Metadata struct {
	NodeName        string `json:"nodeName"`
	NodeDescription string `json:"nodeDescription"`
}

func WellKnownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(WellKnownLinks)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	nodeInfo := NodeInfo{
		Version: "2.0",
		Software: Software{
			Name:    "homeswitch",
			Version: version.Version,
		},
		Protocols: []string{
			"activitypub",
		},
		Services: Services{
			Inbound:  []string{},
			Outbound: []string{},
		},
		OpenRegistration: config.ApprovalRequired,
		Metadata: Metadata{
			NodeName:        config.ServerTitle,
			NodeDescription: config.ServerShortDescription,
		},
	}
	body, err := json.Marshal(nodeInfo)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal NodeInfo")
		http.Error(w, "Failed to marshal NodeInfo", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
