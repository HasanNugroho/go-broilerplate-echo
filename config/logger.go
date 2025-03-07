package config

import (
	"os"

	utils "github.com/HasanNugroho/starter-golang/pkg/utlis"
	"github.com/rs/zerolog"
)

// LoggerConfig ...
type LoggerConfig struct {
	LogLevel string
	Log      *zerolog.Logger
}

func LoadLoggerConfig() (loggerConfig LoggerConfig) {
	loggerConfig.LogLevel = utils.ToString(os.Getenv("LOG_LEVEL"), "error")
	return
}
