package commands

import "testing"

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		name    string
		current string
		latest  string
		want    int
	}{
		{name: "older release", current: "1.0.0-preview", latest: "v1.0.1", want: -1},
		{name: "same release", current: "v1.0.1", latest: "1.0.1", want: 0},
		{name: "newer release", current: "1.0.2", latest: "v1.0.1", want: 1},
		{name: "prerelease lower than release", current: "1.0.0-preview", latest: "1.0.0", want: -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := compareVersions(tc.current, tc.latest)
			if err != nil {
				t.Fatalf("compareVersions returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("compareVersions(%q, %q) = %d, want %d", tc.current, tc.latest, got, tc.want)
			}
		})
	}
}

func TestParseVersionRejectsGarbage(t *testing.T) {
	if _, err := parseVersion("not-a-version"); err == nil {
		t.Fatalf("expected parseVersion to reject garbage")
	}
}
