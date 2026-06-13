package commands

import (
	"fmt"
	"sort"

	"sponsorbucks-client/internal/localconfig"
)

func Status(args []string) {
	cfg, err := localconfig.Load()
	exitOnErr(err)

	fmt.Println("SponsorBucks status")
	fmt.Println("-------------------")
	fmt.Printf("API base URL: %s\n", empty(cfg.APIBaseURL, "not set"))
	fmt.Printf("Device ID:    %s\n", empty(cfg.DeviceID, "not linked"))
	fmt.Printf("User ID:      %s\n", empty(cfg.UserID, "not linked"))
	if cfg.DeviceToken != "" {
		fmt.Println("Auth:         linked")
	} else {
		fmt.Println("Auth:         not linked")
	}
	fmt.Printf("Paused:       %t\n", cfg.Paused)
	fmt.Printf("Disabled:     %s\n", disabledTools(cfg.DisabledTools))
}

func Logout(args []string) {
	cfg, err := localconfig.Load()
	exitOnErr(err)
	cfg.DeviceToken = ""
	cfg.UserID = ""
	cfg.LinkCode = ""
	exitOnErr(localconfig.Save(cfg))
	fmt.Println("Logged out locally. Device registration remains on the server until removed there.")
}

func empty(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func disabledTools(values map[string]bool) string {
	if len(values) == 0 {
		return "none"
	}
	names := make([]string, 0, len(values))
	for name, disabled := range values {
		if disabled {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		return "none"
	}
	return fmt.Sprintf("%v", names)
}
