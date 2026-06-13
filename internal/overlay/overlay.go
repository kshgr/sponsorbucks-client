package overlay

import (
	"fmt"
	"os"
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

func RenderSponsored(line, openURL, openHint string) string {
	line = SponsoredLine(line)
	linkMode := supportsOSC8() && openURL != ""
	budget := sponsoredLineBudget(linkMode, openHint)
	if len([]rune(line)) > budget {
		line = truncateRunes(line, budget)
	}

	rendered := "Sponsored - " + line
	if linkMode {
		rendered += " " + osc8(openURL, "open")
	} else if openHint != "" {
		rendered += " [open: " + openHint + "]"
	}

	if supportsANSIColor() {
		return "\033[32m" + rendered + "\033[0m"
	}
	return rendered
}

func sponsoredLineBudget(linkMode bool, openHint string) int {
	budget := 60 - len([]rune("Sponsored - "))
	if linkMode {
		budget -= len([]rune(" [open]"))
	} else if openHint != "" {
		budget -= len([]rune(" [open: " + openHint + "]"))
	}
	if budget < 0 {
		return 0
	}
	return budget
}

func supportsANSIColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	if fi, err := os.Stderr.Stat(); err == nil && (fi.Mode()&os.ModeCharDevice) != 0 {
		return true
	}
	if fi, err := os.Stdout.Stat(); err == nil && (fi.Mode()&os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

func supportsOSC8() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("WT_SESSION") != "" || os.Getenv("VTE_VERSION") != "" {
		return true
	}
	if strings.Contains(strings.ToLower(os.Getenv("TERM_PROGRAM")), "wezterm") {
		return true
	}
	if strings.Contains(strings.ToLower(os.Getenv("TERM_PROGRAM")), "iterm") {
		return true
	}
	if strings.Contains(strings.ToLower(os.Getenv("TERM_PROGRAM")), "vscode") {
		return true
	}
	if strings.Contains(strings.ToLower(os.Getenv("TERM")), "kitty") {
		return true
	}
	return false
}

func osc8(url, text string) string {
	return "\033]8;;" + url + "\033\\" + text + "\033]8;;\033\\"
}

func truncateRunes(value string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	if max == 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}
