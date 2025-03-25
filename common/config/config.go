package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config holds the application configuration
type Config struct {
	Port               string `envconfig:"PORT" default:"8080"`
	SupabaseURL        string `envconfig:"SUPABASE_URL" required:"true"`
	SupabaseAnonKey    string `envconfig:"SUPABASE_ANON_KEY" required:"true"`
	SupabaseServiceKey string `envconfig:"SUPABASE_SERVICE_ROLE" required:"true"`
	Environment        string `envconfig:"ENV" default:"development"`
	TenantHeader       string `envconfig:"TENANT_HEADER" default:"X-Tenant-ID"`
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
