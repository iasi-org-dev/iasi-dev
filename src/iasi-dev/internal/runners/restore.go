package runners

import (
	"path/filepath"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// Restore restores iasi-org-dev to a previously tagged organization version.
// The implementation is intentionally pending; the CLI contract is established now.
func Restore(Parms *structures.Parms) []string {
	requireTargetVersion(Parms, "restore")
	targets := append([]string{}, Parms.Repos...)

	for _, repository := range targets {
		cli.Header(*Parms, "Restore %s", filepath.Base(repository))
		cli.Step(*Parms, "Restoring %s -> %s", Parms.Version, Parms.TargetVersion)
	}

	cli.Info(*Parms, "Restore %s -> %s pendiente de implementación.", Parms.Version, Parms.TargetVersion)
	RC.Add(Parms.RC, RC.NothingToDo)

	// TODO: verify TargetVersion exists in every repository that belongs to the snapshot.
	// TODO: protect or reject uncommitted local changes before restoring anything.
	// TODO: restore the complete iasi-org-dev state to TargetVersion.
	// TODO: remove repositories from the active workspace when they did not exist in TargetVersion.
	// TODO: never delete tags newer than TargetVersion; historical versions remain immutable.
	// TODO: update the iasi-org-dev VERSION organization variable after a successful restore.

	return targets
}
