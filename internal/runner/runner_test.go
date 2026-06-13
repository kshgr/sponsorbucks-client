package runner

import (
	"testing"
	"time"

	"sponsorbucks-client/internal/adstate"
	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/localconfig"
)

func TestDisabledToolDoesNotActivate(t *testing.T) {
	cfg := localconfig.Config{
		DisabledTools: map[string]bool{
			"claude-code": true,
		},
	}
	if placementEnabled(cfg, "claude-code") {
		t.Fatalf("expected claude-code to be disabled")
	}
	if !placementEnabled(cfg, "codex") {
		t.Fatalf("expected codex to remain enabled")
	}
}

func TestLiveModeRequiresLinkedDevice(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())

	code, err := RunAgent(Options{
		Surface: "codex",
		Command: []string{"cmd", "/c", "exit", "0"},
		Version: buildinfo.Version,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestClickEventUsesCurrentVisibleCreative(t *testing.T) {
	now := time.Date(2026, time.June, 13, 12, 0, 0, 0, time.UTC)
	state := adstate.VisibleAd{
		SessionID:     "sess_123",
		DeviceID:      "dev_123",
		Surface:       "codex",
		CampaignID:    "camp_123",
		CreativeID:    "creative_123",
		CreativeHash:  "hash_123",
		ClientVersion: "1.0.0-preview",
		BuildID:       "dev",
		BuildChannel:  "dev",
	}
	ev := clickEventFromState(state, now)
	if ev.CampaignID != state.CampaignID || ev.CreativeID != state.CreativeID || ev.CreativeHash != state.CreativeHash {
		t.Fatalf("click event did not use the current visible creative: %+v", ev)
	}
	if ev.ClickedAt != now.UTC().Format(time.RFC3339) {
		t.Fatalf("unexpected clicked_at: %s", ev.ClickedAt)
	}
}

func TestNoAdDoesNotCreatePayableImpression(t *testing.T) {
	if hasPayableImpression(noLiveAdFrame()) {
		t.Fatalf("expected no live ad to be non-payable")
	}
}
