package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/runner"
)

func Run(args []string, info buildinfo.Info) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	surface := fs.String("surface", "terminal", "agent/tool surface name")
	demo := fs.Bool("demo", false, "run without sending backend events")
	sticky := fs.String("sticky", "auto", "sticky rendering mode: auto|on|off")
	_ = fs.Parse(args)

	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Println("Usage: sponsorbucks run --surface <surface> -- <agent command>")
		os.Exit(1)
	}

	code, err := runner.RunAgent(runner.Options{
		Surface: *surface,
		Command: rest,
		Version: info.ClientVersion,
		Demo:    *demo,
		Sticky:  parseStickyMode(*sticky),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "sponsorbucks run failed: %v\n", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func parseStickyMode(value string) runner.StickyMode {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto":
		return runner.StickyAuto
	case "on":
		return runner.StickyOn
	case "off":
		return runner.StickyOff
	default:
		fmt.Fprintf(os.Stderr, "invalid --sticky value: %s\n", value)
		os.Exit(1)
		return runner.StickyAuto
	}
}
