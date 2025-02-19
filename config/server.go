package config

import (
	"os"
	"strings"
)

type ServerConfig struct {
	ServerHost     string
	ServerPort     string
	ServerEnv      string
	AllowedOrigins []string
}

// server - port and env
func LoadServerConfig() (serverConfig ServerConfig) {
	serverConfig.ServerHost = strings.TrimSpace(os.Getenv("APP_HOST"))
	serverConfig.ServerPort = strings.TrimSpace(os.Getenv("APP_PORT"))
	serverConfig.ServerEnv = strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	serverConfig.AllowedOrigins = strings.Split(strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS")), ",")
	return
}
