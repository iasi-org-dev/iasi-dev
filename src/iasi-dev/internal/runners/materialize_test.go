package runners

import (
	"path/filepath"
	"testing"
)

func TestMaterializeOrganizationUsesDestinationDirectoryName(t *testing.T) {
	destination := filepath.Join("work", "iasi-org")
	if got := materializeOrganization(destination); got != "iasi-org" {
		t.Fatalf("materializeOrganization() = %q, want iasi-org", got)
	}
}

func TestMaterializeTemporaryIsSiblingOfDestination(t *testing.T) {
	destination := filepath.Join("work", "iasi-org")
	want := filepath.Join("work", ".iasi-org.materialize.tmp")
	if got := materializeTemporary(destination); got != want {
		t.Fatalf("materializeTemporary() = %q, want %q", got, want)
	}
}

func TestMaterializeConfirmed(t *testing.T) {
	accepted := []string{"s", "S", "si", "sí", "y", "YES"}
	for _, value := range accepted {
		if !materializeConfirmed(value) {
			t.Fatalf("materializeConfirmed(%q) = false, want true", value)
		}
	}

	rejected := []string{"", "n", "no", "anything"}
	for _, value := range rejected {
		if materializeConfirmed(value) {
			t.Fatalf("materializeConfirmed(%q) = true, want false", value)
		}
	}
}
