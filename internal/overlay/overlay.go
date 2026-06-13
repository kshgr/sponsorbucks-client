package overlay

import (
	"fmt"
	"strings"
	"time"
)

func SponsoredLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		line = "Sponsored - SponsorBucks founding campaign"
	}
	if len([]rune(line)) > 60 {
		r := []rune(line)
		line = string(r[:60])
	}
	return line
}

func Print(line string) {
	now := time.Now().Format("15:04:05")
	fmt.Printf("\n[%s] %s\n", now, SponsoredLine(line))
}
