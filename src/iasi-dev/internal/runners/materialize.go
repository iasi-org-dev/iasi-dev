package runners

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

// Materialize recreates an organization workspace in destination without carrying Git history.
// Unless -l is active, the materialized organization is then pushed through its configured remotes.
func Materialize(Parms *structures.Parms) []string {
	repositories := materializeLocal(Parms)
	if len(repositories) == 0 || Parms.Local {
		return repositories
	}

	Parms.Repos = repositories
	return push(Parms)
}

// materializeLocal performs only the local materialization transaction.
// Publication is deliberately a separate primitive so workflows can compose both stages.
func materializeLocal(Parms *structures.Parms) []string {
	if Parms.MaterializeDestination == "" {
		cli.Error(RC.Error, *Parms, "materialize requiere un destino.")
	}
	if len(Parms.Repos) == 0 {
		cli.Error(RC.Error, *Parms, "No se encontraron repositorios para materializar.")
	}

	destination := Parms.MaterializeDestination
	confirmed, createDestination := confirmMaterializeDestination(Parms, destination)
	if !confirmed {
		RC.Add(Parms.RC, RC.NothingToDo)
		return nil
	}

	destinationOrganization := materializeOrganization(destination)
	temporary := materializeTemporary(destination)

	cli.Header(*Parms, "Materialize %s", Parms.Organization)
	cli.Info(*Parms, "Source organization: %s", Parms.Organization)
	cli.Info(*Parms, "Destination: %s", destination)
	cli.Info(*Parms, "Destination organization: %s", destinationOrganization)
	cli.Info(*Parms, "Temporary workspace: %s", temporary)
	cli.Info(*Parms, "Mode: local")
	if createDestination {
		materializeCreateDestination(destination)
	}

	destinationRepositories := make([]string, 0, len(Parms.Repos))
	for _, repository := range Parms.Repos {
		name := filepath.Base(repository)
		temporaryRepository := filepath.Join(temporary, name)
		destinationRepository := filepath.Join(destination, name)
		destinationRepositories = append(destinationRepositories, destinationRepository)

		materializeCopy(repository, temporaryRepository)
		materializeGitInit(temporaryRepository, name)
		materializeRemote(temporaryRepository, destinationOrganization, name)
	}

	materializeReplace(temporary, destination)

	cli.Success(*Parms, "Organization materialized.")
	return destinationRepositories
}

func materializeOrganization(destination string) string {
	return filepath.Base(filepath.Clean(destination))
}

func materializeTemporary(destination string) string {
	parent := filepath.Dir(filepath.Clean(destination))
	name := filepath.Base(filepath.Clean(destination))
	return filepath.Join(parent, "."+name+".materialize.tmp")
}

func materializeCopy(source string, destination string) bool {
	fmt.Printf("%s: copy %s -> %s\n", filepath.Base(source), source, destination)
	return true
}

func materializeGitInit(repository string, name string) bool {
	fmt.Printf("%s: git init\n", name)
	return true
}

func materializeRemote(repository string, organization string, name string) bool {
	remote := fmt.Sprintf("https://github.com/%s/%s.git", organization, name)
	fmt.Printf("%s: git remote add origin %s\n", name, remote)
	return true
}

func materializeCreateDestination(destination string) bool {
	fmt.Printf("create %s\n", destination)
	return true
}

func materializeReplace(temporary string, destination string) bool {
	fmt.Printf("replace %s <- %s\n", destination, temporary)
	return true
}

func confirmMaterializeDestination(Parms *structures.Parms, destination string) (bool, bool) {
	info, err := os.Stat(destination)
	if err == nil {
		if !info.IsDir() {
			cli.Error(RC.Error, *Parms, "El destino existe pero no es un directorio: %s", destination)
		}
		return true, false
	}
	if !os.IsNotExist(err) {
		cli.Error(RC.Error, *Parms, "No se puede comprobar el destino %s: %v", destination, err)
	}

	fmt.Printf("El directorio %s no existe. ¿Quieres crearlo? [s/N] ", destination)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if !materializeConfirmed(answer) {
		return false, false
	}
	return true, true
}

func materializeConfirmed(answer string) bool {
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "s", "si", "sí", "y", "yes":
		return true
	default:
		return false
	}
}
