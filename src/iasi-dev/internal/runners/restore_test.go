package runners

import (
	"reflect"
	"testing"

	"iasi-dev/internal/structures"
)

func TestRestoreTargetDefaultsToMain(t *testing.T) {
	Parms := structures.Parms{}
	if got := restoreTarget(&Parms); got != "main" {
		t.Fatalf("restoreTarget() = %q, want main", got)
	}
}

func TestRestoreTargetUsesRequestedVersion(t *testing.T) {
	Parms := structures.Parms{TargetVersion: "v0.5.0"}
	if got := restoreTarget(&Parms); got != "v0.5.0" {
		t.Fatalf("restoreTarget() = %q, want v0.5.0", got)
	}
}

func TestRestoreSwitchArgumentsForMain(t *testing.T) {
	want := []string{"switch", "main"}
	if got := restoreSwitchArguments("main"); !reflect.DeepEqual(got, want) {
		t.Fatalf("restoreSwitchArguments(main) = %v, want %v", got, want)
	}
}

func TestRestoreSwitchArgumentsForTag(t *testing.T) {
	want := []string{"switch", "--detach", "v0.5.0"}
	if got := restoreSwitchArguments("v0.5.0"); !reflect.DeepEqual(got, want) {
		t.Fatalf("restoreSwitchArguments(tag) = %v, want %v", got, want)
	}
}

func TestRestoreRollbackArgumentsForBranch(t *testing.T) {
	state := restoreState{branch: "main"}
	want := []string{"switch", "main"}
	if got := restoreRollbackArguments(state); !reflect.DeepEqual(got, want) {
		t.Fatalf("restoreRollbackArguments(branch) = %v, want %v", got, want)
	}
}

func TestRestoreRollbackArgumentsForDetachedHead(t *testing.T) {
	state := restoreState{commit: "abc123"}
	want := []string{"switch", "--detach", "abc123"}
	if got := restoreRollbackArguments(state); !reflect.DeepEqual(got, want) {
		t.Fatalf("restoreRollbackArguments(detached) = %v, want %v", got, want)
	}
}
