package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port    int    `yaml:"port"`
		Timeout string `yaml:"timeout"`
	} `yaml:"server"`
	MongoDB struct {
		URI      string `yaml:"uri"`
		Database string `yaml:"database"`
	} `yaml:"mongodb"`
	AI struct {
		DefaultModel string `yaml:"default_model"`
		Timeout      string `yaml:"timeout"`
	} `yaml:"ai"`
	Observability struct {
		OTelEnabled  bool   `yaml:"otel_enabled"`
		OTelEndpoint string `yaml:"otel_endpoint"`
		ServiceName  string `yaml:"service_name"`
	} `yaml:"observability"`
}

type Secrets struct {
	GroqAPIKey            string `yaml:"groq_api_key"`
	InfisicalClientID     string `yaml:"infisical_client_id"`
	InfisicalClientSecret string `yaml:"infisical_client_secret"`
}

// LoadConfig reads a YAML config file and applies environment variable
// overrides. Precedence: environment variables > config file.
//
// Supported overrides: SERVER_PORT, MONGO_URI, OTEL_ENABLED,
// OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_SERVICE_NAME.
func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}
	applyEnvOverrides(&cfg)
	return &cfg, nil
}

// LoadSecrets reads a YAML secrets file and applies environment variable
// overrides. A missing file is not an error: secrets may come entirely from
// the environment (12-factor style).
//
// Supported overrides: GROQ_API_KEY, INFISICAL_CLIENT_ID,
// INFISICAL_CLIENT_SECRET.
func LoadSecrets(path string) (*Secrets, error) {
	var sec Secrets
	buf, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(buf, &sec); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	applySecretEnvOverrides(&sec)
	return &sec, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("MONGO_URI"); v != "" {
		cfg.MongoDB.URI = v
	}
	if v, ok := os.LookupEnv("OTEL_ENABLED"); ok {
		if enabled, err := strconv.ParseBool(v); err == nil {
			cfg.Observability.OTelEnabled = enabled
		}
	}
	if v := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); v != "" {
		cfg.Observability.OTelEndpoint = v
	}
	if v := os.Getenv("OTEL_SERVICE_NAME"); v != "" {
		cfg.Observability.ServiceName = v
	}
}

func applySecretEnvOverrides(sec *Secrets) {
	if v := os.Getenv("GROQ_API_KEY"); v != "" {
		sec.GroqAPIKey = v
	}
	if v := os.Getenv("INFISICAL_CLIENT_ID"); v != "" {
		sec.InfisicalClientID = v
	}
	if v := os.Getenv("INFISICAL_CLIENT_SECRET"); v != "" {
		sec.InfisicalClientSecret = v
	}
}
