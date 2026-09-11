package commands

import (
	"fmt"
	"os"
	"os/exec"

	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// RunProtocolLogged executes a command whose exit status already follows the IASI RC protocol.
func RunProtocolLogged(directory string, logFile *os.File, name string, args ...string) structures.Result {
	if debug {
		fmt.Printf("RunProtocolLogged: directory=%s name=%s args=%v\n", directory, name, args)
	}

	writeCommand(logFile, directory, name, args...)

	cmd := exec.Command(name, args...)
	cmd.Dir = directory
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	err := cmd.Run()
	if err == nil {
		return structures.Result{RC: RC.OK}
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		code := exitError.ExitCode()
		if code >= 0 && code <= 0xFF {
			return structures.Result{RC: code}
		}
	}

	return structures.Result{RC: RC.Error}
}
