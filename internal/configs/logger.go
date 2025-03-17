package configs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/HasanNugroho/starter-golang/internal/pkg/hook"
	"github.com/HasanNugroho/starter-golang/internal/pkg/utils"
	"github.com/rs/zerolog"
)

// LoggerConfig ...
type LoggerConfig struct {
	LogLevel string
	Log      *zerolog.Logger
}

var Logger zerolog.Logger

func LoadLoggerConfig() (loggerConfig LoggerConfig) {
	loggerConfig.LogLevel = utils.ToString(os.Getenv("LOG_LEVEL"), "error")
	return
}

func InitLogger(cfg *Configuration) {
	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Logger.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel // Default to INFO if parsing fails
	}
	zerolog.SetGlobalLevel(level) // Set global log level

	// Declare output as io.Writer to support both os.Stdout and ConsoleWriter
	var output io.Writer = os.Stdout

	// Pretty-print logs for non-production environments
	if cfg.AppEnv != "production" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("[%s]", i))
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("| %s |", i)
			},
			PartsExclude: []string{zerolog.TimestampFieldName},
		}
	}

	// Initialize and assign the global logger
	loggerBuilder := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

		// Apply hook **only in production mode**
	if cfg.AppEnv == "production" {
		loggerBuilder = loggerBuilder.Hook(&hook.LoggerHook{})
		fmt.Println("✅ Logger hook enabled in production mode")
	}

	// Assign the global logger
	Logger = loggerBuilder
}
