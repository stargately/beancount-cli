package config

import "os"

const (
	DefaultAPIURL       = "https://beancount.io/api-gateway/"
	DefaultDashboardURL = "https://beancount.io"
)

// Config holds resolved runtime configuration for the CLI.
type Config struct {
	APIURL       string
	DashboardURL string
}

// Load reads configuration from environment variables, falling back to defaults.
//
// Supported overrides:
//
//	BEANCOUNT_API_URL        – GraphQL endpoint (default: https://beancount.io/api-gateway/)
//	BEANCOUNT_DASHBOARD_URL  – Dashboard base URL (default: https://beancount.io)
func Load() *Config {
	apiURL := os.Getenv("BEANCOUNT_API_URL")
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}

	dashboardURL := os.Getenv("BEANCOUNT_DASHBOARD_URL")
	if dashboardURL == "" {
		dashboardURL = DefaultDashboardURL
	}

	return &Config{
		APIURL:       apiURL,
		DashboardURL: dashboardURL,
	}
}
