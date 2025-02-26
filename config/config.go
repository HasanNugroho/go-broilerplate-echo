package config

import (
	"os"
	"strings"

	utils "github.com/HasanNugroho/starter-golang/pkg/utlis"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Version  string
	Database DatabaseConfig
	Server   ServerConfig
	Security SecurityConfig
	AppEnv   string
}

var generalConfig *Configuration

// Load environment variables
func InitEnv() error {
	if os.Getenv("APP_ENV") != "production" {
		return godotenv.Load()
	}
	return nil
}

// InitConfig initializes the application configuration
func InitConfig() (generalConfig *Configuration, err error) {
	if err := InitEnv(); err != nil {
		return nil, err
	}

	config := Configuration{
		Version:  utils.ToString(os.Getenv("VERSION"), "1.0.0"),
		AppEnv:   strings.ToLower(utils.ToString(os.Getenv("APP_ENV"), "development")),
		Server:   LoadServerConfig(),
		Database: loadDatabaseConfig(),
		Security: LoadSecurityConfig(),
	}

	generalConfig = &config
	return generalConfig, nil
}

// GetConfig - return all the config variables
func GetConfig() (cfg *Configuration, err error) {
	if generalConfig != nil {
		return generalConfig, nil
	}

	cfg, err = InitConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
