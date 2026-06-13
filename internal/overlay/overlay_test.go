package overlay

import "testing"

func TestSponsoredLineMaxLength(t *testing.T) {
	got := SponsoredLine("123456789012345678901234567890123456789012345678901234567890EXTRA")
	if len([]rune(got)) > 60 {
		t.Fatalf("expected <= 60 runes, got %d", len([]rune(got)))
	}
}
