package runners

import (
	"fmt"
	"path/filepath"
	"strings"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/commands"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// Publish delegates publication to iasi.quarto once per selected Git repository.
func Publish(Parms *structures.Parms) []string {
	if Parms.Debug {
		fmt.Printf("Publish: repos=%v\n", Parms.Repos)
	}
	targets := []string{}

	for _, repository := range Parms.Repos {
		if isBlackListed(*Parms, repository) {
			continue
		}
		cli.Info(*Parms, "Publicando %s", filepath.Base(repository))

		rc := publishRepository(repository, *Parms)
		if handleRC(Parms, rc) == RC.Skip {
			addToBlackList(Parms, repository)
			continue
		}

		targets = append(targets, repository)
	}

	return targets
}

// publishRepository delegates project discovery, applicability and publication semantics to iasi.quarto.
func publishRepository(repository string, Parms structures.Parms) int {
	parameters := []string{}
	if Parms.Force {
		parameters = append(parameters, "force = TRUE")
	}

	call := "iasi.quarto::publish(" + strings.Join(parameters, ", ") + ")"
	expression := "rc = " + call + "; quit(status = as.integer(rc), save = \"no\")"

	result := commands.RunProtocolLogged(repository, Parms.LogFile, "Rscript", "-e", expression)
	return result.RC
}
