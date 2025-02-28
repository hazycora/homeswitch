package instance_v2

import (
	"fmt"
	"net/http"

	"git.gay/h/homeswitch/internal/config"
	account_model "git.gay/h/homeswitch/internal/models/account"
	"git.gay/h/homeswitch/internal/version"
	"github.com/gin-gonic/gin"
)

type Instance struct {
	Domain        string                `json:"domain"`
	Title         string                `json:"title"`
	Version       string                `json:"version"`
	SourceURL     string                `json:"source_url"`
	Description   string                `json:"description"`
	Usage         InstanceUsage         `json:"usage"`
	Thumbnail     InstanceThumbnail     `json:"thumbnail"`
	Languages     []string              `json:"languages"`
	Configuration InstanceConfiguration `json:"configuration"`
	Registrations InstanceRegistrations `json:"registrations"`
	Contact       InstanceContact       `json:"contact"`
	Rules         []InstanceRule        `json:"rules"`
}

type InstanceUsage struct {
	Users struct {
		ActiveMonth int64 `json:"active_month"`
	} `json:"users"`
}

type InstanceThumbnail struct {
	URL      string `json:"url"`
	Blurhash string `json:"blurhash"`
	Versions struct {
		Standard string `json:"@1x"`
		Large    string `json:"@2x"`
	}
}

type InstanceConfiguration struct {
	URLs struct {
		Streaming string `json:"streaming"`
		Status    string `json:"status"`
	} `json:"urls"`
	Vapid struct {
		PublicKey string `json:"public_key"`
	} `json:"vapid"`
	Accounts struct {
		MaxPinnedStatuses int64 `json:"max_pinned_statuses"`
		MaxFeaturedTags   int64 `json:"max_featured_tags"`
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
	Translation struct {
		Enabled bool `json:"enabled"`
	} `json:"translation"`
	Reactions struct {
		MaxReactions int64 `json:"max_reactions"`
	} `json:"reactions"`
}

type InstanceRegistrations struct {
	Enabled          bool    `json:"enabled"`
	ApprovalRequired bool    `json:"approval_required"`
	Message          string  `json:"message"`
	URL              *string `json:"url"`
}

type InstanceContact struct {
	Email   string                `json:"email"`
	Account account_model.Account `json:"account"`
}

type InstanceRule struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Hint string `json:"hint"`
}

var configuration = InstanceConfiguration{}

func init() {
	configuration.URLs.Streaming = fmt.Sprintf("wss://%s", config.ServerName)
	configuration.Accounts.MaxPinnedStatuses = 32
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
		Domain:      config.ServerName,
		Title:       config.ServerTitle,
		Description: config.ServerDescription,
		Rules:       []InstanceRule{},
		Version:     version.FullVersion,
		Thumbnail: InstanceThumbnail{
			URL: fmt.Sprintf("https://%s/system/static/banner.png", config.ServerName),
		},
		Registrations: InstanceRegistrations{
			Enabled:          config.RegistrationsEnabled,
			ApprovalRequired: config.ApprovalRequired,
		},
		Languages: []string{"en"},
		Contact: InstanceContact{
			Email: config.AdminEmail,
		},
		Configuration: configuration,
	}
	// TODO: proper usage data

	ContactAccount, ok := account_model.GetAccountByUsername(config.AdminUsername)
	if !ok {
		http.Error(c.Writer, "Error getting contact account", http.StatusInternalServerError)
		return
	}
	instance.Contact.Account = *ContactAccount

	c.JSON(http.StatusOK, instance)
}
