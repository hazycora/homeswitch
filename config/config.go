package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ServerName string
)

func init() {
	ServerName = os.Getenv("SERVER_NAME")
}
