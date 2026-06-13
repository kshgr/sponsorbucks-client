package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"sponsorbucks-client/internal/buildinfo"
)

const githubLatestReleaseURL = "https://api.github.com/repos/kshgr/sponsorbucks-client/releases/latest"

type githubRelease struct {
	TagName   string `json:"tag_name"`
	Name      string `json:"name"`
	HTMLURL   string `json:"html_url"`
	Published string `json:"published_at"`
}

type parsedVersion struct {
	major int
	minor int
	patch int
	pre   string
}

var versionPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?$`)

func UpdateCheck(args []string) {
	_ = args
	current := buildinfo.Current().ClientVersion
	release, err := fetchLatestRelease()
	exitOnErr(err)

	cmp, err := compareVersions(current, release.TagName)
	exitOnErr(err)

	fmt.Printf("Current version: %s\n", current)
	fmt.Printf("Latest release:  %s\n", release.TagName)
	if release.HTMLURL != "" {
		fmt.Printf("Release URL:     %s\n", release.HTMLURL)
	}

	switch {
	case cmp < 0:
		fmt.Println("Update available. Manual update only; no auto-update is performed.")
	case cmp == 0:
		fmt.Println("SponsorBucks is up to date.")
	default:
		fmt.Println("Current build is newer than the latest GitHub release.")
	}
}

func fetchLatestRelease() (githubRelease, error) {
	req, err := http.NewRequest(http.MethodGet, githubLatestReleaseURL, nil)
	if err != nil {
		return githubRelease{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "SponsorBucks-Client")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return githubRelease{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return githubRelease{}, fmt.Errorf("latest release request failed: %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return githubRelease{}, err
	}
	if release.TagName == "" {
		return githubRelease{}, fmt.Errorf("latest release response missing tag_name")
	}
	return release, nil
}

func compareVersions(current, latest string) (int, error) {
	cur, err := parseVersion(current)
	if err != nil {
		return 0, err
	}
	l, err := parseVersion(latest)
	if err != nil {
		return 0, err
	}

	if cur.major != l.major {
		return compareInts(cur.major, l.major), nil
	}
	if cur.minor != l.minor {
		return compareInts(cur.minor, l.minor), nil
	}
	if cur.patch != l.patch {
		return compareInts(cur.patch, l.patch), nil
	}
	if cur.pre == l.pre {
		return 0, nil
	}
	if cur.pre == "" {
		return 1, nil
	}
	if l.pre == "" {
		return -1, nil
	}
	return strings.Compare(cur.pre, l.pre), nil
}

func parseVersion(value string) (parsedVersion, error) {
	match := versionPattern.FindStringSubmatch(strings.TrimSpace(value))
	if match == nil {
		return parsedVersion{}, fmt.Errorf("unsupported version format: %q", value)
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return parsedVersion{}, err
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return parsedVersion{}, err
	}
	patch, err := strconv.Atoi(match[3])
	if err != nil {
		return parsedVersion{}, err
	}
	return parsedVersion{
		major: major,
		minor: minor,
		patch: patch,
		pre:   match[4],
	}, nil
}

func compareInts(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
