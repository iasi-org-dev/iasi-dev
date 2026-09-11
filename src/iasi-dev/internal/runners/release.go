package runners

import (
	"fmt"
	"path/filepath"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/commands"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// Release delegates release to iasi.quarto once per selected Git repository.
func Release(Parms *structures.Parms) []string {
	if Parms.Debug {
		fmt.Printf("Release: repos=%v\n", Parms.Repos)
	}
	targets := []string{}

	for _, repository := range Parms.Repos {
		if isBlackListed(*Parms, repository) {
			continue
		}
		cli.Info(*Parms, "Generando release de %s", filepath.Base(repository))

		rc := releaseRepository(repository, *Parms)
		if RC.Has(handleRC(Parms, rc), RC.Skip) {
			addToBlackList(Parms, repository)
			continue
		}

		targets = append(targets, repository)
	}

	return targets
}

// releaseRepository delegates project discovery, applicability and release semantics to iasi.quarto.
func releaseRepository(repository string, Parms structures.Parms) int {
	expression := "rc = iasi.quarto::release(); quit(status = as.integer(rc), save = \"no\")"
	result := commands.RunProtocolLogged(repository, Parms.LogFile, "Rscript", "-e", expression)
	return result.RC
}
