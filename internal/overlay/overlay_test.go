package overlay

import (
	"strings"
	"testing"
)

func TestSponsoredLineMaxLength(t *testing.T) {
	got := SponsoredLine("123456789012345678901234567890123456789012345678901234567890EXTRA")
	if len([]rune(got)) > 60 {
		t.Fatalf("expected <= 60 runes, got %d", len([]rune(got)))
	}
}

func TestRenderSponsoredNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	got := RenderSponsored("Sponsored - demo creative", "", "sponsorbucks open")
	if strings.Contains(got, "\x1b[32m") {
		t.Fatalf("expected ANSI green to be disabled under NO_COLOR: %q", got)
	}
	if strings.Contains(got, "\x1b]8;;") {
		t.Fatalf("expected OSC8 hyperlink to be disabled under NO_COLOR: %q", got)
	}
}
