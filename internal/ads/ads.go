package ads

import (
	"sponsorbucks-client/internal/hashutil"
	"sponsorbucks-client/internal/overlay"
)

func CreativeHash(line string) string {
	return hashutil.SHA256Hex(overlay.SponsoredLine(line))
}
