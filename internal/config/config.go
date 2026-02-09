package config

import (
	"os"
	"time"
)

const (
	DefaultAPIBase = "https://api.github.com"
	DefaultTimeout = 15 * time.Second
)

// Config holds resolved runtime configuration.
type Config struct {
	Token   string
	APIBase string
	Timeout time.Duration
	JSON    bool
	Verbose bool
}

// ResolveToken returns the token from the first available source:
// flag value > GITHUB_TOKEN > GH_TOKEN.
// Returns empty string if none are set.
func ResolveToken(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		return v
	}
	if v := os.Getenv("GH_TOKEN"); v != "" {
		return v
	}
	return ""
}
