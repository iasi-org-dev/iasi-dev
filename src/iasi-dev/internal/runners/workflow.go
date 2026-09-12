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
		case "promote":
			workflowPromote(true, Parms)
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

// workflowPromote promotes locally and will later publish the resulting iasi-org.
func workflowPromote(standalone bool, Parms *structures.Parms) {
	Parms.Repos = Promote(Parms)
	if len(Parms.Repos) == 0 {
		return
	}

	// TODO: add the atomic operation that publishes the local iasi-org to GitHub.
	// TODO: compose that operation here after Promote succeeds.
	if standalone {
		cli.Step(*Parms, "Publishing promoted organization pending")
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
	case "promote":
		return "Promoting"
	default:
		return name
	}
}
