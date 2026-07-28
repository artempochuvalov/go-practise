package config

import (
	"strings"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://url")
	t.Setenv("PORT", "8000")

	config, err := Load()
	if err != nil {
		t.Fatalf("expected no err, received: %v", err)
	}

	if config.Port != ":8000" {
		t.Fatalf("expected port: %s, received %s", ":8000", config.Port)
	}

	if config.DatabaseURL != "postgresql://url" {
		t.Fatalf("expected db url: %s, received %s", "postgresql://url", config.DatabaseURL)
	}
}

func TestLoadConfig_EmptyUrl(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("PORT", "8000")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no db url provided") {
		t.Fatalf("expected err with message: %s, received: %s", "no db url provided", err.Error())
	}
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://url")
	t.Setenv("PORT", "abc")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected err, received nil")
	}
}

func TestLoadConfig_DefaultPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://url")

	config, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.Port != ":8080" {
		t.Fatalf("expected port: %s, received %s", ":8080", config.Port)
	}
}
