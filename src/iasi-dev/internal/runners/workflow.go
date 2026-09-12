package runners

import (
	"path/filepath"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// Workflow executes the selected workflow one Git repository at a time.
func Workflow(Parms *structures.Parms) {
	repositories := append([]string{}, Parms.Repos...)

	for _, repository := range repositories {
		Parms.Repos = []string{repository}
		cli.Header(*Parms, "%s %s", workflowName(Parms.Subcommand), filepath.Base(repository))

		switch Parms.Subcommand {
		case "build":
			workflowBuild(true, Parms)
		case "publish":
			workflowPublish(true, Parms)
		case "release":
			workflowRelease(true, Parms)
		default:
			cli.Error(RC.Error, *Parms, "Workflow desconocido: %q", Parms.Subcommand)
		}
	}
}

// workflowBuild builds and commits when standalone or used as a checkpoint.
func workflowBuild(standalone bool, Parms *structures.Parms) {
	Parms.LastRC = RC.OK
	Parms.Repos = Build(Parms)
	if !workflowContinuesAfterBuild(Parms.LastRC) {
		Parms.Repos = nil
		return
	}
	if standalone || Parms.Checkpoints {
		Parms.Repos = Commit(Parms)
	}
}

// workflowPublish optionally builds first, publishes and commits when required.
func workflowPublish(standalone bool, Parms *structures.Parms) {
	if Parms.All {
		workflowBuild(false, Parms)
	}
	if len(Parms.Repos) == 0 {
		return
	}

	Parms.Repos = Publish(Parms)
	if standalone || Parms.Checkpoints {
		Parms.Repos = Commit(Parms)
	}
}

// workflowRelease optionally runs previous stages, releases and commits when required.
func workflowRelease(standalone bool, Parms *structures.Parms) {
	if Parms.All {
		workflowPublish(false, Parms)
	}
	if len(Parms.Repos) == 0 {
		return
	}

	Parms.Repos = Release(Parms)
	if standalone || Parms.Checkpoints {
		Parms.Repos = Commit(Parms)
	}
}

// workflowContinuesAfterBuild reports whether later workflow stages apply.
// NothingToDo means build succeeded but found no buildable IASI project.
func workflowContinuesAfterBuild(rc int) bool {
	return RC.Result(rc) != RC.NothingToDo
}

// workflowName returns the display name of a workflow.
func workflowName(name string) string {
	switch name {
	case "build":
		return "Building"
	case "publish":
		return "Publishing"
	case "release":
		return "Releasing"
	default:
		return name
	}
}
