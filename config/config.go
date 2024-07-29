package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	DBConnectionUri        = os.Getenv("DB_CONN_URI")
	ServerName             = os.Getenv("SERVER_NAME")
	ServerURL              = fmt.Sprintf("https://%s", ServerName)
	ServerTitle            = os.Getenv("SERVER_TITLE")
	AdminEmail             = os.Getenv("ADMIN_EMAIL")
	ServerDescription      = os.Getenv("SERVER_DESCRIPTION")
	ServerShortDescription = os.Getenv("SERVER_SHORT_DESCRIPTION")
	AdminUsername          = os.Getenv("ADMIN_USERNAME")

	RegistrationsEnabled = os.Getenv("REGISTRATIONS_ENABLED") == "true"
	ApprovalRequired     = os.Getenv("APPROVAL_REQUIRED") == "true"
	InvitesEnabled       = os.Getenv("INVITES_ENABLED") == "true"
)
