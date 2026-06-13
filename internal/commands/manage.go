package commands

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"sponsorbucks-client/internal/api"
	"sponsorbucks-client/internal/buildinfo"
	"sponsorbucks-client/internal/daemon"
	"sponsorbucks-client/internal/localconfig"
	"sponsorbucks-client/internal/logs"
	sbtools "sponsorbucks-client/internal/tools"
)

func Install(args []string) {
	fs := parseFlagSet("install")
	dryRun := fs.Bool("dry-run", false, "recommended first step: show intended changes without writing")
	yes := fs.Bool("yes", false, "skip PATH block confirmation")
	_ = fs.Parse(args)

	cfgDir, err := localconfig.ConfigDir()
	exitOnErr(err)
	shimDir := filepath.Join(cfgDir, "shims")
	profilePath, shell := sbtools.DetectShellProfile()
	detected := sbtools.DetectInstalled()
	sponsorbucksPath, err := sbtools.SponsorBucksExecutable()
	exitOnErr(err)
	if sbtools.WouldSelfShadow(shimDir, detected) {
		exitOnErr(fmt.Errorf("refusing to install because an installed tool already resolves inside %s", shimDir))
	}

	blockExists, err := profileHasBlock(profilePath)
	exitOnErr(err)
	fmt.Println("Detected tools:")
	printToolPaths(detected)

	files := plannedInstallFiles(profilePath, detected, shimDir, !blockExists)
	printFilesChanged("Files that would change:", files)
	if *dryRun {
		fmt.Println("Dry run only. No files were changed.")
		_ = logs.Append("install", map[string]string{
			"dry_run": "true",
			"files":   fmt.Sprintf("%d", len(files)),
		})
		return
	}

	if err := os.MkdirAll(cfgDir, 0700); err != nil {
		exitOnErr(err)
	}
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		exitOnErr(err)
	}

	changed, err := sbtools.CreateShims(shimDir, sponsorbucksPath, detected, false)
	exitOnErr(err)
	files = changed

	if !blockExists {
		fmt.Printf("Warning: SponsorBucks will edit your shell profile: %s\n", profilePath)
		if *yes {
			exitOnErr(sbtools.AddPathBlock(profilePath, shell, shimDir, false))
			files = append(files, profilePath)
		} else if interactiveInput() && confirm(fmt.Sprintf("Add SponsorBucks PATH block to %s? [y/N] ", profilePath)) {
			exitOnErr(sbtools.AddPathBlock(profilePath, shell, shimDir, false))
			files = append(files, profilePath)
		} else {
			fmt.Println("Skipped shell PATH block.")
		}
	}

	printFilesChanged("Changed files:", files)
	_ = logs.Append("install", map[string]string{
		"dry_run": "false",
		"files":   fmt.Sprintf("%d", len(files)),
	})
}

func Uninstall(args []string) {
	fs := parseFlagSet("uninstall")
	deleteConfig := fs.Bool("delete-config", false, "delete local config as well")
	_ = fs.Parse(args)

	cfgDir, err := localconfig.ConfigDir()
	exitOnErr(err)
	shimDir := filepath.Join(cfgDir, "shims")
	profilePath, _ := sbtools.DetectShellProfile()

	removed := make([]string, 0, len(sbtools.Supported)+2)

	files, err := sbtools.RemoveShims(shimDir, false)
	exitOnErr(err)
	removed = append(removed, files...)

	changed, err := sbtools.RemovePathBlock(profilePath, false)
	exitOnErr(err)
	if changed {
		removed = append(removed, profilePath)
	}

	if *deleteConfig {
		configPath, err := localconfig.ConfigPath()
		exitOnErr(err)
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			exitOnErr(err)
		}
		removed = append(removed, configPath)
	}

	pidPath := filepath.Join(cfgDir, "daemon.pid")
	if err := os.Remove(pidPath); err == nil {
		removed = append(removed, pidPath)
	}

	printFilesChanged("Removed files:", removed)
	_ = logs.Append("uninstall", map[string]string{
		"deleted_config": fmt.Sprintf("%t", *deleteConfig),
		"files":          fmt.Sprintf("%d", len(removed)),
	})
}

func Tools(args []string) {
	_ = args
	cfg, err := localconfig.Load()
	exitOnErr(err)
	detected := sbtools.DetectInstalled()
	shimDir := shimDirFromConfig()
	fmt.Println("SponsorBucks tools")
	fmt.Println("------------------")
	fmt.Printf("Paused: %t\n", cfg.Paused)
	for _, name := range sbtools.Supported {
		enabled := !cfg.DisabledTools[sbtools.CanonicalSurface(name)]
		_, installed := detected[name]
		shimExists := false
		if shimDir != "" {
			shim := sbtools.ShimFilePath(shimDir, name)
			_, shimErr := os.Stat(shim)
			shimExists = shimErr == nil
		}
		fmt.Printf("%s: detected=%t enabled=%t shim=%t path=%s\n", name, installed, enabled, shimExists, detectedOrDash(detected[name]))
	}
}

func EnableTool(args []string) {
	setToolEnabled(args, true, "enable")
}

func DisableTool(args []string) {
	setToolEnabled(args, false, "disable")
}

func Pause(args []string) {
	_ = args
	cfg, err := localconfig.Load()
	exitOnErr(err)
	cfg.Paused = true
	exitOnErr(localconfig.Save(cfg))
	fmt.Println("SponsorBucks paused.")
	_ = logs.Append("pause", nil)
}

func Resume(args []string) {
	_ = args
	cfg, err := localconfig.Load()
	exitOnErr(err)
	cfg.Paused = false
	exitOnErr(localconfig.Save(cfg))
	fmt.Println("SponsorBucks resumed.")
	_ = logs.Append("resume", nil)
}

func Daemon(args []string) {
	fs := parseFlagSet("daemon")
	addr := fs.String("addr", "127.0.0.1:18181", "listen address")
	_ = fs.Parse(args)

	exitOnErr(daemon.Run(daemon.Options{Addr: *addr}))
}

func Privacy(args []string) {
	_ = args
	fmt.Println("SponsorBucks never collects:")
	fmt.Println("- code")
	fmt.Println("- prompts")
	fmt.Println("- model responses")
	fmt.Println("- terminal output")
	fmt.Println("- repo names")
	fmt.Println("- filenames")
	fmt.Println("- screenshots")
	fmt.Println("- clipboard contents")
	fmt.Println("- window titles")
	fmt.Println("- environment variables")
	fmt.Println("- secrets")
}

func Logs(args []string) {
	_ = args
	content, err := logs.Read()
	exitOnErr(err)
	if strings.TrimSpace(content) == "" {
		fmt.Println("No local SponsorBucks logs.")
		return
	}
	fmt.Print(content)
}

func Doctor(args []string) {
	_ = args
	cfg, cfgErr := localconfig.Load()
	detected := sbtools.DetectInstalled()

	ok := true
	configPath, _ := localconfig.ConfigPath()
	report := func(label string, pass bool, detail string) {
		status := "ok"
		if !pass {
			status = "fail"
			ok = false
		}
		if detail != "" {
			fmt.Printf("%s: %s (%s)\n", label, status, detail)
			return
		}
		fmt.Printf("%s: %s\n", label, status)
	}

	_, err := os.Stat(configPath)
	report("config exists", err == nil, configPath)
	if cfgErr != nil {
		report("config load", false, cfgErr.Error())
		cfg = localconfig.Config{}
	}
	report("device linked", cfg.DeviceID != "" && cfg.DeviceToken != "" && cfg.DevicePublicKey != "", cfg.DeviceID)
	report("build type", true, buildinfo.BuildChannel)

	if cfg.APIBaseURL == "" {
		report("api reachable", false, "api base url not set")
	} else {
		client := api.New(cfg.APIBaseURL, cfg.DeviceToken)
		err = client.Health()
		report("api reachable", err == nil, cfg.APIBaseURL)
		if cfg.DeviceToken != "" && cfg.DeviceID != "" && cfg.DevicePrivateKey != "" {
			ad, adErr := client.NextAd(cfg.DeviceID, "codex", "")
			report("ads-next reachable", adErr == nil, ad.CampaignID)
			report("click endpoint reachable", client.PathReachable("/events-click") == nil, "events-click")
		} else {
			report("ads-next reachable", false, "device not linked")
			report("click endpoint reachable", false, "device not linked")
		}
	}

	report("daemon running", daemonHealthy("127.0.0.1:18181"), "127.0.0.1:18181")
	report("supported tools detected", len(detected) > 0, detectedSummary(detected))
	report("shims installed", shimsInstalled(detected), filepath.Join(cfgDirFromConfig(), "shims"))
	report("shell integration active", shellIntegrationActive(), shellProfilePath())
	report("latest release reachable", latestReleaseReachable() == nil, "kshgr/sponsorbucks-client")

	if !ok {
		os.Exit(1)
	}
}

func setToolEnabled(args []string, enabled bool, command string) {
	if len(args) != 1 {
		fmt.Printf("Usage: sponsorbucks %s <tool>\n", command)
		os.Exit(1)
	}
	tool := args[0]
	if !sbtools.IsSupported(tool) {
		fmt.Printf("Unsupported tool: %s\n", tool)
		os.Exit(1)
	}
	cfg, err := localconfig.Load()
	exitOnErr(err)
	if cfg.DisabledTools == nil {
		cfg.DisabledTools = make(map[string]bool)
	}
	canonical := sbtools.CanonicalSurface(tool)
	cfg.DisabledTools[canonical] = !enabled
	exitOnErr(localconfig.Save(cfg))
	_ = logs.Append(command, map[string]string{
		"enabled": fmt.Sprintf("%t", enabled),
		"tool":    canonical,
	})
	if enabled {
		fmt.Printf("Enabled %s.\n", tool)
		return
	}
	fmt.Printf("Disabled %s.\n", tool)
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

func interactiveInput() bool {
	fi, err := os.Stdin.Stat()
	return err == nil && (fi.Mode()&os.ModeCharDevice) != 0
}

func parseFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	return fs
}

func printFilesChanged(header string, files []string) {
	seen := make(map[string]bool)
	unique := make([]string, 0, len(files))
	for _, file := range files {
		if file == "" || seen[file] {
			continue
		}
		seen[file] = true
		unique = append(unique, file)
	}
	fmt.Println(header)
	if len(unique) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, file := range unique {
		fmt.Printf("  %s\n", file)
	}
}

func printToolPaths(detected map[string]string) {
	if len(detected) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, name := range sbtools.Supported {
		if path, ok := detected[name]; ok {
			fmt.Printf("  %s -> %s\n", name, path)
		}
	}
}

func profileHasBlock(profilePath string) (bool, error) {
	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return strings.Contains(string(data), "SponsorBucks PATH START"), nil
}

func cfgDirFromConfig() string {
	dir, err := localconfig.ConfigDir()
	if err != nil {
		return ""
	}
	return dir
}

func shimDirFromConfig() string {
	dir := cfgDirFromConfig()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "shims")
}

func shellProfilePath() string {
	path, _ := sbtools.DetectShellProfile()
	return path
}

func daemonHealthy(addr string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + addr + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 300
}

func shimsInstalled(detected map[string]string) bool {
	dir := shimDirFromConfig()
	if dir == "" {
		return false
	}
	for name := range detected {
		if _, err := os.Stat(sbtools.ShimFilePath(dir, name)); err != nil {
			return false
		}
	}
	return len(detected) > 0
}

func shellIntegrationActive() bool {
	path, _ := sbtools.DetectShellProfile()
	ok, err := profileHasBlock(path)
	return err == nil && ok
}

func plannedInstallFiles(profilePath string, detected map[string]string, shimDir string, includeProfile bool) []string {
	files := make([]string, 0, len(detected)+1)
	for _, tool := range sbtools.Supported {
		if _, ok := detected[tool]; ok {
			files = append(files, sbtools.ShimFilePath(shimDir, tool))
		}
	}
	if includeProfile {
		files = append(files, profilePath)
	}
	return files
}

func detectedOrDash(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

func detectedSummary(values map[string]string) string {
	if len(values) == 0 {
		return "none"
	}
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func latestReleaseReachable() error {
	_, err := fetchLatestRelease()
	return err
}
