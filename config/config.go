package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ServerName             = os.Getenv("SERVER_NAME")
	ServerTitle            = os.Getenv("SERVER_TITLE")
	AdminEmail             = os.Getenv("ADMIN_EMAIL")
	ServerDescription      = os.Getenv("SERVER_DESCRIPTION")
	ServerShortDescription = os.Getenv("SERVER_SHORT_DESCRIPTION")
	AdminUsername          = os.Getenv("ADMIN_USERNAME")
)
