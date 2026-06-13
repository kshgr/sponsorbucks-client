package system

// Signals is intentionally privacy-preserving. Do not add fields that expose
// code, prompts, terminal output, filenames, repository names, screenshots,
// clipboard contents, or window titles.
type Signals struct {
	ScreenUnlocked             bool
	RecentInputBucket          string
	ForegroundSupportedSurface bool
}

// CollectSignals returns conservative v0 attention signals.
//
// Because `sponsorbucks run` owns the terminal session, ForegroundSupportedSurface
// can be true for wrapper-based placements. More precise OS-specific checks
// should be added later behind small platform files.
func CollectSignals() Signals {
	return Signals{
		ScreenUnlocked:             true,
		RecentInputBucket:          "unknown_v0",
		ForegroundSupportedSurface: true,
	}
}
