package db

import "testing"

func TestBuildDSN(t *testing.T) {
	cfg := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "secret",
		Name:     "guitar_specs",
		SSLMode:  "disable",
	}

	dsn := buildDSN(cfg)
	expected := "postgres://postgres:secret@localhost:5432/guitar_specs?sslmode=disable"
	if dsn != expected {
		t.Fatalf("expected %s, got %s", expected, dsn)
	}
}
