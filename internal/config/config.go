// Package config loads runtime settings from environment variables. Secrets
// (source API keys/tokens) come ONLY from the environment — never hardcoded and
// never sent to the frontend.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all runtime settings.
type Config struct {
	Addr            string        // listen address, e.g. ":8080"
	SourceTimeout   time.Duration // per-source deadline during fan-out
	DefaultCurrency string
	CORSOrigin      string // allowed browser origin for the React app

	// Source credentials — populated in later phases. Kept server-side only.
	KiwiAPIKey         string
	TravelpayoutsToken string

	EnableMock bool // serve deterministic offers when no real source is configured
}

// Load reads configuration from the environment, applying sane defaults.
func Load() Config {
	return Config{
		Addr:               env("FM_ADDR", ":8080"),
		SourceTimeout:      envDuration("FM_SOURCE_TIMEOUT", 8*time.Second),
		DefaultCurrency:    env("FM_DEFAULT_CURRENCY", "RUB"),
		CORSOrigin:         env("FM_CORS_ORIGIN", "http://localhost:5173"),
		KiwiAPIKey:         os.Getenv("FM_KIWI_API_KEY"),
		TravelpayoutsToken: os.Getenv("FM_TRAVELPAYOUTS_TOKEN"),
		EnableMock:         envBool("FM_ENABLE_MOCK", true),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
