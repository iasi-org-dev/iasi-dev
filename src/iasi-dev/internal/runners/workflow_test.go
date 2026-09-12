package runners

import (
	"path/filepath"
	"testing"

	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

func TestWorkflowStopsAfterBuildNothingToDo(t *testing.T) {
	if workflowContinuesAfterBuild(RC.NothingToDo) {
		t.Fatal("workflow must stop after build returns NothingToDo")
	}
}

func TestWorkflowContinuesAfterOtherBuildResults(t *testing.T) {
	for _, rc := range []int{RC.OK, RC.Info, RC.Warning, RC.Attention} {
		if !workflowContinuesAfterBuild(rc) {
			t.Fatalf("workflow unexpectedly stops after build RC 0x%X", rc)
		}
	}
}

func TestWorkflowOrganizationRoot(t *testing.T) {
	repositories := []string{
		filepath.Join("C:", "iasi-org-dev", "repo-a"),
		filepath.Join("C:", "iasi-org-dev", "repo-b"),
	}

	got := workflowOrganizationRoot(repositories)
	want := filepath.Join("C:", "iasi-org-dev")
	if got != want {
		t.Fatalf("workflowOrganizationRoot() = %q, want %q", got, want)
	}
}

func TestWorkflowPromoteDestination(t *testing.T) {
	rc := RC.OK
	parms := structures.Parms{
		Organization: "iasi-org-dev",
		Repos: []string{
			filepath.Join("C:", "iasi-org-dev", "repo-a"),
			filepath.Join("C:", "iasi-org-dev", "repo-b"),
		},
		RC: &rc,
	}

	got := workflowPromoteDestination(&parms)
	want := filepath.Join("C:", "iasi-org")
	if got != want {
		t.Fatalf("workflowPromoteDestination() = %q, want %q", got, want)
	}
}
