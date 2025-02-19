package utils

import (
	"strconv"
	"strings"
	"time"
)

// ToInt converts a string to an int with a default value.
func ToInt(value string, defaultValue int) int {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return defaultValue
	}
	return v
}

// ToBool converts a string to a bool with a default value.
func ToBool(value string, defaultValue bool) bool {
	v, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return defaultValue
	}
	return v
}

// ToFloat converts a string to a float64 with a default value.
func ToFloat(value string, defaultValue float64) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToDuration converts a string to a time.Duration with a default value.
func ToDuration(value string, defaultValue time.Duration) time.Duration {
	v, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return defaultValue
	}
	return v
}

// ToString trims a string and returns a default value if empty.
func ToString(value string, defaultValue string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return defaultValue
	}
	return v
}
