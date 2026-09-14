package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigAppliesEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
server:
  port: 9000
mongodb:
  uri: "mongodb://file-host:27017"
  database: "lumina_plane"
ai:
  default_model: "llama3-8b-8192"
observability:
  otel_enabled: false
  otel_endpoint: "localhost:4317"
  service_name: "lumina-plane-server"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("SERVER_PORT", "8123")
	t.Setenv("MONGO_URI", "mongodb://env-host:27017")
	t.Setenv("OTEL_ENABLED", "true")

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if cfg.Server.Port != 8123 {
		t.Errorf("Server.Port = %d, want 8123 (env override)", cfg.Server.Port)
	}
	if cfg.MongoDB.URI != "mongodb://env-host:27017" {
		t.Errorf("MongoDB.URI = %q, want env override", cfg.MongoDB.URI)
	}
	if !cfg.Observability.OTelEnabled {
		t.Error("Observability.OTelEnabled = false, want true (env override)")
	}
	if cfg.Observability.OTelEndpoint != "localhost:4317" {
		t.Errorf("Observability.OTelEndpoint = %q, want file value kept when no env override", cfg.Observability.OTelEndpoint)
	}
}

func TestLoadConfigWithoutEnvKeepsFileValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := `
server:
  port: 9000
mongodb:
  uri: "mongodb://file-host:27017"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if cfg.Server.Port != 9000 {
		t.Errorf("Server.Port = %d, want 9000 from file", cfg.Server.Port)
	}
	if cfg.MongoDB.URI != "mongodb://file-host:27017" {
		t.Errorf("MongoDB.URI = %q, want file value", cfg.MongoDB.URI)
	}
}

func TestLoadSecretsToleratesMissingFileAndAppliesEnv(t *testing.T) {
	t.Setenv("GROQ_API_KEY", "gsk_test_123")

	sec, err := LoadSecrets(filepath.Join(t.TempDir(), "does_not_exist.yaml"))
	if err != nil {
		t.Fatalf("LoadSecrets should tolerate a missing file: %v", err)
	}
	if sec.GroqAPIKey != "gsk_test_123" {
		t.Errorf("GroqAPIKey = %q, want env override", sec.GroqAPIKey)
	}
}

func TestLoadSecretsReadsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secrets.yaml")
	if err := os.WriteFile(path, []byte("groq_api_key: \"file_key\"\n"), 0o600); err != nil {
		t.Fatalf("write secrets: %v", err)
	}

	sec, err := LoadSecrets(path)
	if err != nil {
		t.Fatalf("LoadSecrets: %v", err)
	}
	if sec.GroqAPIKey != "file_key" {
		t.Errorf("GroqAPIKey = %q, want file value", sec.GroqAPIKey)
	}
}
