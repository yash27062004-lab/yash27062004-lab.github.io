package config

import "testing"

func TestDefaultConfigIsValid(t *testing.T) {
	cfg := Default()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if cfg.OutputDir != "dist" {
		t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, "dist")
	}

	if cfg.SiteURL != "https://yash27062004-lab.github.io" {
		t.Errorf("unexpected default site URL: %q", cfg.SiteURL)
	}
}

func TestLoadUsesDefaultsWhenEnvironmentIsEmpty(t *testing.T) {
	keys := []string{
		"SITE_URL",
		"SITE_OUTPUT_DIR",
		"SITE_TEMPLATES_DIR",
		"SITE_CONTENT_DIR",
		"SITE_STATIC_DIR",
	}

	for _, key := range keys {
		t.Setenv(key, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg != Default() {
		t.Errorf("Load() = %#v, want %#v", cfg, Default())
	}
}

func TestLoadAppliesEnvironmentOverrides(t *testing.T) {
	t.Setenv("SITE_URL", "https://example.test/")
	t.Setenv("SITE_OUTPUT_DIR", "build")
	t.Setenv("SITE_TEMPLATES_DIR", "views")
	t.Setenv("SITE_CONTENT_DIR", "articles")
	t.Setenv("SITE_STATIC_DIR", "assets")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.SiteURL != "https://example.test" {
		t.Errorf("SiteURL = %q", cfg.SiteURL)
	}

	if cfg.OutputDir != "build" {
		t.Errorf("OutputDir = %q", cfg.OutputDir)
	}

	if cfg.TemplatesDir != "views" {
		t.Errorf("TemplatesDir = %q", cfg.TemplatesDir)
	}

	if cfg.ContentDir != "articles" {
		t.Errorf("ContentDir = %q", cfg.ContentDir)
	}

	if cfg.StaticDir != "assets" {
		t.Errorf("StaticDir = %q", cfg.StaticDir)
	}
}

func TestValidateRejectsEmptySiteURL(t *testing.T) {
	cfg := Default()
	cfg.SiteURL = ""

	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() expected an error for an empty site URL")
	}
}

func TestLoadRejectsInvalidSiteURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "missing scheme",
			url:  "localhost:8080",
		},
		{
			name: "HTTP instead of HTTPS",
			url:  "http://example.test",
		},
		{
			name: "missing host",
			url:  "https://",
		},
		{
			name: "query string",
			url:  "https://example.test/path?preview=true",
		},
		{
			name: "fragment",
			url:  "https://example.test/#section",
		},
		{
			name: "credentials",
			url:  "https://user:pass@example.test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("SITE_URL", test.url)

			_, err := Load()
			if err == nil {
				t.Fatalf("Load() with URL %q should fail", test.url)
			}
		})
	}
}
