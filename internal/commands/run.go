package commands

import (
	"flag"
	"fmt"
	"os"

	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/runner"
)

func Run(args []string, info buildinfo.Info) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	surface := fs.String("surface", "terminal", "agent/tool surface name")
	demo := fs.Bool("demo", false, "run without sending backend events")
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
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "sponsorbucks run failed: %v\n", err)
		os.Exit(1)
	}
	os.Exit(code)
}
