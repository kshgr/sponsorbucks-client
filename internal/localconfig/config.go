package localconfig

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	APIBaseURL       string          `json:"api_base_url"`
	DeviceID         string          `json:"device_id"`
	UserID           string          `json:"user_id"`
	DeviceToken      string          `json:"device_token"`
	DevicePublicKey  string          `json:"device_public_key"`
	DevicePrivateKey string          `json:"device_private_key"`
	LinkCode         string          `json:"link_code"`
	Paused           bool            `json:"paused"`
	DisabledTools    map[string]bool `json:"disabled_tools,omitempty"`
}

func (c Config) Redacted() Config {
	out := c
	out.DeviceToken = redactValue(out.DeviceToken)
	out.DevicePrivateKey = redactValue(out.DevicePrivateKey)
	out.LinkCode = redactValue(out.LinkCode)
	return out
}

func redactValue(value string) string {
	if value == "" {
		return ""
	}
	return "***redacted***"
}

func Load() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Save(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sponsorbucks"), nil
}

func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}
