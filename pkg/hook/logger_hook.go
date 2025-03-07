package hook

import (
	"fmt"

	"github.com/rs/zerolog"
)

type LoggerHook struct{}

func (t LoggerHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if level > zerolog.WarnLevel {
		// Modify event immediately
		e.Str("severity", level.String())

		// Send notification asynchronously
		go notify(level.String(), message)
	}
}

func notify(title, msg string) error {
	fmt.Println(title, "hook triggered:", msg)
	return nil
}
