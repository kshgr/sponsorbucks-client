package ads

import (
	"testing"

	"sponsorbucks-client/internal/hashutil"
	"sponsorbucks-client/internal/overlay"
)

func TestCreativeHashMatchesVisibleLine(t *testing.T) {
	line := "Sponsored - demo creative"
	want := hashutil.SHA256Hex(overlay.SponsoredLine(line))
	got := CreativeHash(line)
	if got != want {
		t.Fatalf("creative hash mismatch: got %s want %s", got, want)
	}
}
