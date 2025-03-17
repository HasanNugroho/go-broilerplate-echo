package configs

import (
	"os"
	"strings"

	"github.com/HasanNugroho/starter-golang/internal/pkg/utils"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Version  string
	Database RDBMSConfig
	Redis    RedisConfig
	Server   ServerConfig
	Security SecurityConfig
	AppEnv   string
	Logger   LoggerConfig
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
		Database: loadRDBMSConfig(),
		Redis:    loadRedisConfig(),
		Security: LoadSecurityConfig(),
		Logger:   LoadLoggerConfig(),
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
