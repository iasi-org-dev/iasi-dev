package runners

import (
	"testing"

	"iasi-dev/internal/consts/RC"
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
