package tools

import "strings"

var surfaceAliases = map[string]string{
	"claude":      "claude-code",
	"claude-code": "claude-code",
	"gemini":      "gemini-cli",
	"gemini-cli":  "gemini-cli",
	"codex":       "codex",
	"pi":          "pi",
	"aider":       "aider",
	"opencode":    "opencode",
	"terminal":    "generic-terminal",
}

var commandAliases = map[string]string{
	"claude-code": "claude",
	"gemini-cli":  "gemini",
}

func CanonicalSurface(surface string) string {
	surface = strings.ToLower(strings.TrimSpace(surface))
	if canonical, ok := surfaceAliases[surface]; ok {
		return canonical
	}
	return surface
}

func CommandForSurface(surface string) string {
	surface = CanonicalSurface(surface)
	if command, ok := commandAliases[surface]; ok {
		return command
	}
	return surface
}

func SurfaceAliases() map[string]string {
	out := make(map[string]string, len(surfaceAliases))
	for k, v := range surfaceAliases {
		out[k] = v
	}
	return out
}
