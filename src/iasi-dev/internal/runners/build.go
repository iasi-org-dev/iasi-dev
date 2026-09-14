package runners

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/commands"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// Build delegates build to iasi.quarto once per selected Git repository.
func Build(Parms *structures.Parms) []string {
	if Parms.Debug {
		fmt.Printf("Build: repos=%v\n", Parms.Repos)
	}
	targets := []string{}

	for _, repository := range Parms.Repos {
		if isBlackListed(*Parms, repository) {
			continue
		}
		if Parms.Subcommand == "" {
			cli.Header(*Parms, "Build %s", filepath.Base(repository))
		}
		cli.Step(*Parms, "Building")

		rc := buildRepository(repository, *Parms)
		Parms.LastRC = rc
		handled := handleRC(Parms, rc)
		if RC.Has(handled, RC.Skip) {
			addToBlackList(Parms, repository)
			continue
		}
		targets = append(targets, repository)
	}

	return targets
}

// buildRepository delegates project discovery, types and build semantics to iasi.quarto.
func buildRepository(repository string, Parms structures.Parms) int {
	parameters := []string{}

	if Parms.Format != "" {
		formats := strings.Split(Parms.Format, ",")
		for i, format := range formats {
			formats[i] = strconv.Quote(strings.TrimSpace(format))
		}
		parameters = append(parameters, "format = c("+strings.Join(formats, ", ")+")")
	}

	call := "iasi.quarto::build(" + strings.Join(parameters, ", ") + ")"
	expression := "rc = " + call + "; quit(status = as.integer(rc), save = \"no\")"

	result := commands.RunLogged(repository, Parms.LogFile, "Rscript", "-e", expression)
	return result.RC
}
