package runners

import (
	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// Workflow executes the selected workflow one Git repository at a time.
func Workflow(Parms *structures.Parms) {
	repositories := append([]string{}, Parms.Repos...)

	for _, repository := range repositories {
		Parms.Repos = []string{repository}

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
	Parms.Repos = Build(Parms)
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
