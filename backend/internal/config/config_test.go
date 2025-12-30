package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Env != "development" {
		t.Fatalf("expected default env development, got %s", cfg.Env)
	}
	if cfg.ServerPort != "8080" {
		t.Fatalf("expected default port 8080, got %s", cfg.ServerPort)
	}
	if cfg.DBHost != "localhost" {
		t.Fatalf("expected default db host localhost, got %s", cfg.DBHost)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("GUITAR_SPECS_ENV", "test")
	t.Setenv("GUITAR_SPECS_SERVER_PORT", "9999")
	t.Setenv("GUITAR_SPECS_DB_HOST", "db")
	t.Setenv("GUITAR_SPECS_DB_NAME", "guitars")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Env != "test" {
		t.Fatalf("expected env test, got %s", cfg.Env)
	}
	if cfg.ServerPort != "9999" {
		t.Fatalf("expected port 9999, got %s", cfg.ServerPort)
	}
	if cfg.DBHost != "db" {
		t.Fatalf("expected db host db, got %s", cfg.DBHost)
	}
	if cfg.DBName != "guitars" {
		t.Fatalf("expected db name guitars, got %s", cfg.DBName)
	}
}
