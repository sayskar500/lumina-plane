package config

import (
	"os"

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
}

type Secrets struct {
	GroqAPIKey          string `yaml:"groq_api_key"`
	InfisicalClientID     string `yaml:"infisical_client_id"`
	InfisicalClientSecret string `yaml:"infisical_client_secret"`
}

func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadSecrets(path string) (*Secrets, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sec Secrets
	err = yaml.Unmarshal(buf, &sec)
	if err != nil {
		return nil, err
	}
	return &sec, nil
}
