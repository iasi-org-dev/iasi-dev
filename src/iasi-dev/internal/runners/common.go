package runners

import (
	"fmt"
	"path/filepath"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// handleRC accumulates one result and applies the common error policy.
func handleRC(Parms *structures.Parms, rc int) int {
	if Parms.Debug {
		fmt.Printf("handleRC: tolerant=%t rc=0x%02X\n", Parms.Tolerant, rc)
	}

	RC.Add(Parms.RC, rc)

	if !RC.IsErroneous(rc) {
		return rc
	}

	if Parms.Tolerant {
		return RC.Skip
	}

	cli.ErrorMessage(*Parms, "La operación ha fallado con RC 0x%02X. Revisa el log: %s", rc, logName(*Parms))
	return RC.Skip
}

// logName returns the current log path when available.
func logName(Parms structures.Parms) string {
	if Parms.LogFile == nil {
		return "(sin log)"
	}
	return Parms.LogFile.Name()
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
