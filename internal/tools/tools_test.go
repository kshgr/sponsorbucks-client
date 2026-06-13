package tools

import (
	"path/filepath"
	"testing"
)

func TestWouldSelfShadow(t *testing.T) {
	shimDir := filepath.Join(t.TempDir(), "shims")
	detected := map[string]string{
		"codex": filepath.Join(shimDir, "codex"),
	}
	if !WouldSelfShadow(shimDir, detected) {
		t.Fatalf("expected self-shadow detection to trigger")
	}
}
