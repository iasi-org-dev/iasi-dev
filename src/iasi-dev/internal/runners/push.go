package runners

import (
	"fmt"
	"path/filepath"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// push publishes the current local organization state through each repository's origin remote.
// It is deliberately package-private: there is no public `iasi-dev push` command.
// Mutating Git operations are still previewed, consistently with promote and materialize.
func push(Parms *structures.Parms) []string {
	if len(Parms.Repos) == 0 {
		cli.Error(RC.Error, *Parms, "No se encontraron repositorios para publicar.")
	}

	cli.Step(*Parms, "Pushing organization")
	for _, repository := range Parms.Repos {
		if !pushRepository(repository) {
			cli.Error(RC.Error, *Parms, "No se pudo publicar %s.", repository)
		}
	}

	cli.Success(*Parms, "Organization pushed.")
	return append([]string{}, Parms.Repos...)
}

func pushRepository(repository string) bool {
	fmt.Printf("%s: git push -u origin HEAD\n", filepath.Base(repository))
	return true
}
