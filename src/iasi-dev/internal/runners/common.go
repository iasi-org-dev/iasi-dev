package runners

import (
	"fmt"
	"path/filepath"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// checkTolerant stops the execution when tolerance is disabled or skips the current item.
func checkTolerant(Parms structures.Parms, rc int) int {
	if Parms.Debug {
		fmt.Printf("checkTolerant: tolerant=%t rc=%d\n", Parms.Tolerant, rc)
	}
	if Parms.Tolerant {
		return RC.Skip
	}

	cli.Error(rc, Parms, "La operación ha fallado.")
	return RC.Skip
}

// addToBlackList blacklists repository.
func addToBlackList(Parms *structures.Parms, repository string) {
	if repository == "" || isBlackListed(*Parms, repository) {
		return
	}
	Parms.BlackList = append(Parms.BlackList, filepath.Clean(repository))
}

// isBlackListed reports whether repository has been invalidated by a previous operation.
func isBlackListed(Parms structures.Parms, repository string) bool {
	for _, blocked := range Parms.BlackList {
		if filepath.Clean(blocked) == filepath.Clean(repository) {
			return true
		}
	}
	return false
}
