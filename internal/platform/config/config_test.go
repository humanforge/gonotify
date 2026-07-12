package config

import (
	"os"
	"testing"
)

func TestLoad_RequiredVars(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "test")
	os.Setenv("DATABASE_CONN_URL", "postgres://localhost:5432/test")
	os.Setenv("LOG_ROOT_PATH", "/tmp/logs")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DatabaseURL != "postgres://localhost:5432/test" {
		t.Errorf("expected DatabaseURL 'postgres://localhost:5432/test', got %q", cfg.DatabaseURL)
	}
	if cfg.LogRootPath != "/tmp/logs" {
		t.Errorf("expected LogRootPath '/tmp/logs', got %q", cfg.LogRootPath)
	}
	if cfg.Port != ":8080" {
		t.Errorf("expected Port ':8080', got %q", cfg.Port)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DATABASE_CONN_URL")
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "test")
	os.Setenv("DATABASE_CONN_URL", "postgres://localhost:5432/test")
	os.Setenv("LOG_ROOT_PATH", "/tmp/logs")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.LogRootPath != "/tmp/logs" {
		t.Errorf("expected LogRootPath '/tmp/logs', got %q", cfg.LogRootPath)
	}
	if cfg.ShutdownTimeout == 0 {
		t.Error("expected non-zero ShutdownTimeout")
	}
}

func TestLoad_PortFormatting(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "test")
	os.Setenv("DATABASE_CONN_URL", "postgres://localhost:5432/test")
	os.Setenv("LOG_ROOT_PATH", "/tmp/logs")
	os.Setenv("PORT", "3000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Port != ":3000" {
		t.Errorf("expected Port ':3000', got %q", cfg.Port)
	}
}

func TestValidLogLevel(t *testing.T) {
	tests := []struct {
		level string
		want  bool
	}{
		{"debug", true},
		{"info", true},
		{"warn", true},
		{"error", true},
		{"fatal", true},
		{"unknown", false},
		{"", false},
	}
	for _, tt := range tests {
		got := ValidLogLevel(tt.level)
		if got != tt.want {
			t.Errorf("ValidLogLevel(%q) = %v, want %v", tt.level, got, tt.want)
		}
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_INT", "42")
	if got := getEnvInt("TEST_INT", 0); got != 42 {
		t.Errorf("getEnvInt(TEST_INT) = %d, want 42", got)
	}
	if got := getEnvInt("MISSING", 99); got != 99 {
		t.Errorf("getEnvInt(MISSING) = %d, want 99", got)
	}
	os.Setenv("TEST_BAD", "abc")
	if got := getEnvInt("TEST_BAD", 10); got != 10 {
		t.Errorf("getEnvInt(TEST_BAD) = %d, want 10", got)
	}
}
