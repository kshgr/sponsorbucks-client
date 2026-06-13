package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"sponsorbucks-client/internal/localconfig"
	"sponsorbucks-client/internal/logs"
)

func Config(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sponsorbucks config set-api <url> | sponsorbucks config show")
		os.Exit(1)
	}

	switch args[0] {
	case "set-api":
		if len(args) != 2 {
			fmt.Println("Usage: sponsorbucks config set-api <url>")
			os.Exit(1)
		}
		cfg, err := localconfig.Load()
		exitOnErr(err)
		cfg.APIBaseURL = args[1]
		exitOnErr(localconfig.Save(cfg))
		fmt.Println("API base URL saved.")
		_ = logs.Append("config-set-api", nil)
	case "show":
		if len(args) != 1 {
			fmt.Println("Usage: sponsorbucks config show")
			os.Exit(1)
		}
		cfg, err := localconfig.Load()
		exitOnErr(err)
		out, err := json.MarshalIndent(cfg.Redacted(), "", "  ")
		exitOnErr(err)
		fmt.Println(string(out))
	default:
		fmt.Println("Unknown config command.")
		os.Exit(1)
	}
}
