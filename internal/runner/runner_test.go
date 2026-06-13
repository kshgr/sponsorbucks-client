package runner

import (
	"testing"

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
