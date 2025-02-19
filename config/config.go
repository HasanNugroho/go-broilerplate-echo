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

var configAll *Configuration

// Load environment variables
func InitEnv() error {
	return godotenv.Load()
}

// InitConfig initializes the application configuration
func InitConfig() error {
	if err := InitEnv(); err != nil {
		return err
	}

	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return err
	}

	config := Configuration{
		Version:  utils.ToString(os.Getenv("VERSION"), "1.0.0"),
		AppEnv:   strings.ToLower(utils.ToString(os.Getenv("APP_ENV"), "development")),
		Server:   LoadServerConfig(),
		Database: dbConfig,
		Security: LoadSecurityConfig(),
	}

	configAll = &config
	return nil
}

// GetConfig - return all the config variables
func GetConfig() *Configuration {
	return configAll
}
