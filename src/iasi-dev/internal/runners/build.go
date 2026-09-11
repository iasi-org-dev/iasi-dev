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
		cli.Info(*Parms, "Construyendo %s", filepath.Base(repository))

		rc := buildRepository(repository, *Parms)
		if rc == RC.Skip {
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

	expression := "iasi.quarto::build(" + strings.Join(parameters, ", ") + ")"
	result := commands.RunFriendlyLogged(repository, Parms.LogFile, "Rscript", "-e", expression)
	if RC.IsErroneous(result.RC) {
		return checkTolerant(Parms, RC.Build)
	}

	return RC.OK
}
