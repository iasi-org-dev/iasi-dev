package runners

import "testing"

func TestParseSemanticVersion(t *testing.T) {
	tests := []struct {
		value string
		ok    bool
	}{
		{"v0.5.0", true},
		{"v1.0.12", true},
		{"v10.20.30", true},
		{"0.5.0", false},
		{"v0.5", false},
		{"v0.5.0.1", false},
		{"v01.5.0", false},
		{"v0.05.0", false},
		{"v0.5.00", false},
		{"v0.x.0", false},
	}

	for _, test := range tests {
		_, ok := parseSemanticVersion(test.value)
		if ok != test.ok {
			t.Fatalf("parseSemanticVersion(%q) ok = %t, want %t", test.value, ok, test.ok)
		}
	}
}

func TestCompareSemanticVersions(t *testing.T) {
	v030, _ := parseSemanticVersion("v0.3.0")
	v060, _ := parseSemanticVersion("v0.6.0")
	v070, _ := parseSemanticVersion("v0.7.0")
	v100, _ := parseSemanticVersion("v1.0.0")

	if compareSemanticVersions(v030, v060) >= 0 {
		t.Fatal("v0.3.0 must be lower than v0.6.0")
	}
	if compareSemanticVersions(v060, v060) != 0 {
		t.Fatal("v0.6.0 must equal v0.6.0")
	}
	if compareSemanticVersions(v070, v060) <= 0 {
		t.Fatal("v0.7.0 must be greater than v0.6.0")
	}
	if compareSemanticVersions(v100, v070) <= 0 {
		t.Fatal("v1.0.0 must be greater than v0.7.0")
	}
}

func TestHighestSemanticVersion(t *testing.T) {
	fallback, _ := parseSemanticVersion("v0.3.0")
	highest, text := highestSemanticVersion([]string{"other", "v0.1.0", "v0.6.0", "v0.4.0"}, fallback, "v0.3.0")
	want, _ := parseSemanticVersion("v0.6.0")

	if compareSemanticVersions(highest, want) != 0 {
		t.Fatal("highest semantic version must be v0.6.0")
	}
	if text != "v0.6.0" {
		t.Fatalf("highest text = %q, want v0.6.0", text)
	}
}

func TestHighestSemanticVersionFallsBackToCurrentVersion(t *testing.T) {
	fallback, _ := parseSemanticVersion("v0.5.0")
	highest, text := highestSemanticVersion([]string{"not-a-version", "release"}, fallback, "v0.5.0")

	if compareSemanticVersions(highest, fallback) != 0 {
		t.Fatal("highest semantic version must fall back to current version")
	}
	if text != "v0.5.0" {
		t.Fatalf("highest text = %q, want v0.5.0", text)
	}
}
