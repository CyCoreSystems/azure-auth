package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
	"gopkg.in/yaml.v3"
)

// Config is the configuration for the client OAUTH2 system
type Config struct {
	Username     string          `yaml:"username"`
	TenantID     string          `yaml:"tenantID"`
	ClientID     string          `yaml:"clientID"`
	ClientSecret string          `yaml:"clientSecret"`
	Scopes       []string        `yaml:"scopes"`
	Redirect     *RedirectConfig `yaml:"redirect"`
}

// RedirectConfig describes the OAUTH2 delegationr redirect setup (from client config on Microsoft)
type RedirectConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// URL returns the RedirectURL.
func (rc *RedirectConfig) URL() string {
	return fmt.Sprintf("http://%s:%d/%s", rc.Host, rc.Port, strings.Trim(rc.Path, "/"))
}

// OAuth2 returns the OAuth2 config from this configuration.
func (cfg *Config) OAuth2() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		RedirectURL:  cfg.Redirect.URL(),
		Endpoint:     microsoft.AzureADEndpoint(cfg.TenantID),
	}
}

// LoadConfig loads the configuration from the default configuration file.
func LoadConfig() (*Config, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determind configuration directory")
	}

	fn := path.Join(cfgDir, "azure", "config.yaml")

	f, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file %q: %w", fn, err)
	}

	cfg := new(Config)

	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %q: %w", fn, err)
	}

	return cfg, nil
}
