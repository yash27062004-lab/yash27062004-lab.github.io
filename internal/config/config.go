package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	OutputDir    string
	TemplatesDir string
	ContentDir   string
	StaticDir    string
	SiteURL      string
}

func Default() Config {
	return Config{
		OutputDir:    "dist",
		TemplatesDir: "templates",
		ContentDir:   "content",
		StaticDir:    "static",
		SiteURL:      "https://yash27062004-lab.github.io",
	}
}

func Load() (Config, error) {
	cfg := Default()

	cfg.OutputDir = envOrDefault("SITE_OUTPUT_DIR", cfg.OutputDir)
	cfg.TemplatesDir = envOrDefault("SITE_TEMPLATES_DIR", cfg.TemplatesDir)
	cfg.ContentDir = envOrDefault("SITE_CONTENT_DIR", cfg.ContentDir)
	cfg.StaticDir = envOrDefault("SITE_STATIC_DIR", cfg.StaticDir)
	cfg.SiteURL = envOrDefault("SITE_URL", cfg.SiteURL)

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	cfg.SiteURL = strings.TrimRight(cfg.SiteURL, "/")

	return cfg, nil
}

func (c Config) Validate() error {
	fields := []struct {
		name  string
		value string
	}{
		{"output directory", c.OutputDir},
		{"templates directory", c.TemplatesDir},
		{"content directory", c.ContentDir},
		{"static directory", c.StaticDir},
		{"site URL", c.SiteURL},
	}

	for _, field := range fields {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("%s cannot be empty", field.name)
		}
	}

	parsed, err := url.Parse(strings.TrimSpace(c.SiteURL))
	if err != nil {
		return fmt.Errorf("invalid site URL: %w", err)
	}

	if parsed.Scheme != "https" || parsed.Host == "" {
		return fmt.Errorf("site URL must be an absolute HTTPS URL")
	}

	if parsed.User != nil ||
		parsed.RawQuery != "" ||
		parsed.Fragment != "" {
		return fmt.Errorf("site URL cannot contain credentials, a query, or a fragment")
	}

	return nil
}

func envOrDefault(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}

	return value
}
