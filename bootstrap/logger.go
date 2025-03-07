package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/pkg/hook"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger(cfg *config.Configuration) {
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
	Logger = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger().
		Hook(&hook.LoggerHook{})
}
