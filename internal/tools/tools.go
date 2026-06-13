package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

var Supported = []string{"codex", "claude", "pi", "aider", "opencode", "gemini"}

const shimMarker = "SponsorBucks managed shim"

func IsSupported(name string) bool {
	switch CanonicalSurface(name) {
	case "codex", "claude-code", "gemini-cli", "pi", "aider", "opencode", "generic-terminal":
		return true
	default:
		return false
	}
}

func DetectInstalled() map[string]string {
	found := make(map[string]string, len(Supported))
	for _, name := range Supported {
		path, err := exec.LookPath(name)
		if err == nil && path != "" {
			found[name] = path
		}
	}
	return found
}

func SortedNames(values map[string]string) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func ShimFilePath(shimDir, tool string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(shimDir, tool+".cmd")
	}
	return filepath.Join(shimDir, tool)
}

func SponsorBucksExecutable() (string, error) {
	return os.Executable()
}

func RemoveShims(shimDir string, dryRun bool) ([]string, error) {
	var removed []string
	for _, tool := range Supported {
		path := ShimFilePath(shimDir, tool)
		existing, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return removed, err
		}
		if !strings.Contains(string(existing), shimMarker) {
			continue
		}
		if !dryRun {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return removed, err
			}
		}
		removed = append(removed, path)
	}
	return removed, nil
}

func WouldSelfShadow(shimDir string, detected map[string]string) bool {
	if shimDir == "" {
		return false
	}
	shimDir = normalizePath(shimDir)
	for _, path := range detected {
		if path == "" {
			continue
		}
		if pathWithin(normalizePath(path), shimDir) {
			return true
		}
	}
	return false
}

func DetectShellProfile() (string, string) {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		if strings.Contains(strings.ToLower(os.Getenv("PSModulePath")), "powershell") {
			return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"), "powershell"
		}
		return filepath.Join(home, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"), "powershell"
	}
	shell := filepath.Base(os.Getenv("SHELL"))
	switch shell {
	case "zsh":
		return filepath.Join(home, ".zshrc"), "zsh"
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish"), "fish"
	default:
		return filepath.Join(home, ".bashrc"), "bash"
	}
}

func PathBlock(shell string, shimDir string) string {
	markerStart := "# >>> SponsorBucks PATH START >>>"
	markerEnd := "# <<< SponsorBucks PATH END <<<"
	switch shell {
	case "powershell":
		return fmt.Sprintf("%s\n$env:Path = \"%s;$env:Path\"\n%s\n", markerStart, escapePSPath(shimDir), markerEnd)
	case "fish":
		return fmt.Sprintf("%s\nfish_add_path %q\n%s\n", markerStart, shimDir, markerEnd)
	default:
		return fmt.Sprintf("%s\nexport PATH=%q:$PATH\n%s\n", markerStart, shimDir, markerEnd)
	}
}

func AddPathBlock(profilePath, shell, shimDir string, dryRun bool) error {
	if dryRun {
		return nil
	}
	existing, err := os.ReadFile(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if strings.Contains(string(existing), "SponsorBucks PATH START") {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(profilePath), 0755); err != nil {
		return err
	}
	updated := append(existing, []byte("\n"+PathBlock(shell, shimDir))...)
	return os.WriteFile(profilePath, updated, 0644)
}

func RemovePathBlock(profilePath string, dryRun bool) (bool, error) {
	existing, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	updated := removeBlock(string(existing))
	if updated == string(existing) {
		return false, nil
	}
	if dryRun {
		return true, nil
	}
	return true, os.WriteFile(profilePath, []byte(updated), 0644)
}

func removeBlock(content string) string {
	start := strings.Index(content, "# >>> SponsorBucks PATH START >>>")
	end := strings.Index(content, "# <<< SponsorBucks PATH END <<<")
	if start == -1 || end == -1 || end < start {
		return content
	}
	end += len("# <<< SponsorBucks PATH END <<<")
	updated := content[:start] + content[end:]
	return strings.TrimSpace(updated) + "\n"
}

func shimContent(tool, sponsorbucksPath, realPath string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("@echo off\r\nREM %s\r\nsetlocal\r\n\"%s\" run --surface %s -- \"%s\" %%*\r\n", shimMarker, sponsorbucksPath, tool, realPath)
	}
	return fmt.Sprintf("#!/bin/sh\n# %s\nexec %q run --surface %s -- %q \"$@\"\n", shimMarker, sponsorbucksPath, tool, realPath)
}

func escapePSPath(path string) string {
	return strings.ReplaceAll(path, "`", "``")
}

func CreateShims(shimDir, sponsorbucksPath string, detected map[string]string, dryRun bool) ([]string, error) {
	paths := make([]string, 0, len(Supported))
	for _, tool := range Supported {
		_, ok := detected[tool]
		if !ok {
			continue
		}
		path := ShimFilePath(shimDir, tool)
		if !dryRun {
			if err := ensureManagedWritable(path); err != nil {
				return nil, err
			}
		}
		paths = append(paths, path)
	}
	if dryRun {
		return paths, nil
	}
	for _, tool := range Supported {
		realPath, ok := detected[tool]
		if !ok {
			continue
		}
		path := ShimFilePath(shimDir, tool)
		content := shimContent(tool, sponsorbucksPath, realPath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return paths, err
		}
		if err := os.WriteFile(path, []byte(content), 0755); err != nil {
			return paths, err
		}
	}
	return paths, nil
}

func ensureManagedWritable(path string) error {
	existing, err := os.ReadFile(path)
	if err == nil {
		if !strings.Contains(string(existing), shimMarker) {
			return fmt.Errorf("refusing to overwrite non-SponsorBucks file: %s", path)
		}
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func normalizePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(abs)
}

func pathWithin(path, parent string) bool {
	if path == parent {
		return true
	}
	sep := string(os.PathSeparator)
	parent = strings.TrimRight(parent, sep)
	return strings.HasPrefix(path, parent+sep)
}
