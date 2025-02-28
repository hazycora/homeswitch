package instance_v1

import (
	"fmt"
	"net/http"

	"git.gay/h/homeswitch/internal/config"
	account_model "git.gay/h/homeswitch/internal/models/account"
	"git.gay/h/homeswitch/internal/version"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Instance struct {
	URI              string `json:"uri"`
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	Email            string `json:"email"`
	Version          string `json:"version"`
	URLs             struct {
		StreamingAPI string `json:"streaming_api"`
	} `json:"urls"`
	Stats struct {
		UserCount   int64 `json:"user_count"`
		StatusCount int64 `json:"status_count"`
		DomainCount int64 `json:"domain_count"`
	} `json:"stats"`
	Thumbnail            string                `json:"thumbnail"`
	Languages            []string              `json:"languages"`
	RegistrationsEnabled bool                  `json:"registrations"`
	ApprovalRequired     bool                  `json:"approval_required"`
	InvitesEnabled       bool                  `json:"invites_enabled"`
	Configuration        InstanceConfiguration `json:"configuration"`
	ContactAccount       account_model.Account `json:"contact_account"`
	Rules                []InstanceRule        `json:"rules"`
}

type InstanceConfiguration struct {
	Accounts struct {
		MaxFeaturedTags int64 `json:"max_featured_tags"`
	} `json:"accounts"`
	Statuses struct {
		MaxCharacters            int64    `json:"max_characters"`
		MaxMediaAttachments      int64    `json:"max_media_attachments"`
		CharactersReservedPerURL int64    `json:"characters_reserved_per_url"`
		SupportedMimeTypes       []string `json:"supported_mime_types"`
	} `json:"statuses"`
	MediaAttachments struct {
		SupportedMimeTypes  []string `json:"supported_mime_types"`
		ImageSizeLimit      int64    `json:"image_size_limit"`
		ImageMatrixLimit    int64    `json:"image_matrix_limit"`
		VideoSizeLimit      int64    `json:"video_size_limit"`
		VideoFrameRateLimit int64    `json:"video_frame_rate_limit"`
		VideoMatrixLimit    int64    `json:"video_matrix_limit"`
	} `json:"media_attachments"`
	Polls struct {
		AllowMedia             bool  `json:"allow_media"`
		MaxOptions             int64 `json:"max_options"`
		MaxCharactersPerOption int64 `json:"max_characters_per_option"`
		MinExpiration          int64 `json:"min_expiration"`
		MaxExpiration          int64 `json:"max_expiration"`
	} `json:"polls"`
	Reactions struct {
		MaxReactions int64 `json:"max_reactions"`
	} `json:"reactions"`
}

type InstanceRule struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Hint string `json:"hint"`
}

var configuration = InstanceConfiguration{}

func init() {
	configuration.Accounts.MaxFeaturedTags = 4
	configuration.Statuses.MaxCharacters = 4096
	configuration.Statuses.MaxMediaAttachments = 4
	configuration.Statuses.CharactersReservedPerURL = 23
	configuration.Statuses.SupportedMimeTypes = []string{"text/plain", "text/markdown", "text/html"}
	configuration.MediaAttachments.SupportedMimeTypes = []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"video/mp4",
		"video/webm",
		"video/quicktime",
		"video/ogg",
		"audio/webm",
		"audio/wave",
		"audio/wav",
		"audio/x-wav",
		"audio/x-pn-wave",
		"audio/vnd.wave",
		"audio/ogg",
		"audio/vorbis",
		"audio/mpeg",
		"audio/mp3",
		"audio/flac",
		"audio/aac",
		"audio/m4a",
		"audio/x-m4a",
		"audio/mp4",
		"audio/3gpp",
		"audio/x-ms-asf",
		"audio/opus",
	}
	configuration.MediaAttachments.ImageSizeLimit = 8 * 1024 * 1024
	configuration.MediaAttachments.ImageMatrixLimit = 33177600 // 7680x4320px, taken from masto
	configuration.MediaAttachments.VideoSizeLimit = 40 * 1024 * 1024
	configuration.MediaAttachments.VideoFrameRateLimit = 120
	configuration.MediaAttachments.VideoMatrixLimit = 8294400 // 3840x2160px, taken from masto
	configuration.Polls.AllowMedia = true
	configuration.Polls.MaxOptions = 4
	configuration.Polls.MaxCharactersPerOption = 100
	configuration.Polls.MinExpiration = 300     // 5 minutes
	configuration.Polls.MaxExpiration = 2592000 // 30 days
	configuration.Reactions.MaxReactions = 4
}

func Handler(c *gin.Context) {
	instance := Instance{
		URI:                  config.ServerName,
		Title:                config.ServerTitle,
		ShortDescription:     config.ServerShortDescription,
		Description:          config.ServerDescription,
		Rules:                []InstanceRule{},
		Version:              version.FullVersion,
		Thumbnail:            fmt.Sprintf("https://%s/system/static/banner.png", config.ServerName),
		Languages:            []string{"en"},
		RegistrationsEnabled: config.RegistrationsEnabled,
		ApprovalRequired:     config.ApprovalRequired,
		InvitesEnabled:       config.InvitesEnabled,
		Email:                config.AdminEmail,
		Configuration:        configuration,
	}
	instance.URLs.StreamingAPI = fmt.Sprintf("wss://%s", config.ServerName)

	userCount, err := account_model.GetLocalAccountCount()
	if err != nil {
		log.Error().Err(err).Msg("Could not get local account count")
		http.Error(c.Writer, "Error getting user count", http.StatusInternalServerError)
		return
	}
	instance.Stats.UserCount = userCount
	// TODO: status count, domain count

	ContactAccount, ok := account_model.GetAccountByUsername(config.AdminUsername)
	if !ok {
		http.Error(c.Writer, "Error getting contact account", http.StatusInternalServerError)
		return
	}
	instance.ContactAccount = *ContactAccount

	c.JSON(http.StatusOK, instance)
}
