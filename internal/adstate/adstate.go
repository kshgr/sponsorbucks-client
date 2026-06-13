package adstate

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"sponsorbucks-client/internal/localconfig"
)

type VisibleAd struct {
	Visible       bool   `json:"visible"`
	SessionID     string `json:"session_id,omitempty"`
	DeviceID      string `json:"device_id,omitempty"`
	Surface       string `json:"surface,omitempty"`
	CampaignID    string `json:"campaign_id,omitempty"`
	CreativeID    string `json:"creative_id,omitempty"`
	CreativeHash  string `json:"creative_hash,omitempty"`
	Line          string `json:"line,omitempty"`
	OpenURL       string `json:"open_url,omitempty"`
	UpdatedAt     string `json:"updated_at"`
	ClientVersion string `json:"client_version,omitempty"`
	BuildID       string `json:"build_id,omitempty"`
	BuildChannel  string `json:"build_channel,omitempty"`
}

func Path() (string, error) {
	dir, err := localconfig.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "ad_state.json"), nil
}

func Load() (VisibleAd, error) {
	path, err := Path()
	if err != nil {
		return VisibleAd{}, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return VisibleAd{}, nil
	}
	if err != nil {
		return VisibleAd{}, err
	}
	var state VisibleAd
	if err := json.Unmarshal(data, &state); err != nil {
		return VisibleAd{}, err
	}
	return state, nil
}

func Save(state VisibleAd) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if state.UpdatedAt == "" {
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
