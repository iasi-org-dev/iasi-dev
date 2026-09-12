package args

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"iasi-dev/internal/cli"
	"iasi-dev/internal/consts"
	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
	"iasi-dev/internal/tools"
)

// Parse parses command-line arguments. Preparation is performed later by main.
func Parse(command string, values []string) structures.Parms {
	subcommand := ""

	if command == "workflow" {
		if len(values) == 0 {
			rc := RC.OK
			cli.Error(RC.Error, structures.Parms{RC: &rc}, "Falta el comando del workflow.")
		}

		subcommand = values[0]
		values = values[1:]
	}

	Parms := parseArguments(values)
	Parms.Subcommand = subcommand
	extractTargetVersion(command, &Parms)
	extractVersionOrganization(command, &Parms)
	Parms.RequestedTargets = append([]string{}, Parms.Targets...)

	return Parms
}

// extractTargetVersion separates organization version operands from filesystem targets.
func extractTargetVersion(command string, Parms *structures.Parms) {
	requiresVersion := command == "promote" || command == "restore" || (command == "workflow" && Parms.Subcommand == "promote")
	if !requiresVersion || len(Parms.Targets) == 0 {
		return
	}

	Parms.TargetVersion = Parms.Targets[0]
	Parms.Targets = Parms.Targets[1:]
}

// extractVersionOrganization separates the optional organization operand from filesystem targets.
func extractVersionOrganization(command string, Parms *structures.Parms) {
	if command != "version" || len(Parms.Targets) == 0 {
		return
	}
	if len(Parms.Targets) > 1 {
		cli.Error(RC.Error, *Parms, "version acepta como máximo una organización.")
	}

	Parms.Organization = Parms.Targets[0]
	Parms.Targets = nil
}

func parseArguments(args []string) structures.Parms {
	rc := RC.OK
	Parms := structures.Parms{
		Verbose:    1,
		RC:         &rc,
		Exclusions: append([]string{}, consts.RequiredExclusions...),
	}

	for i := 0; i < len(args); i++ {
		if len(args[i]) == 0 {
			invalidArgument(&Parms, args[i])
		}

		switch args[i][0] {
		case '-':
			parseFlagOrParameter(args, &i, &Parms)
		default:
			parseTarget(args, i, &Parms)
		}
	}

	return Parms
}

func parseTarget(args []string, i int, Parms *structures.Parms) {
	Parms.Targets = append(Parms.Targets, args[i])
}

func parseFlagOrParameter(args []string, i *int, Parms *structures.Parms) {
	switch len(args[*i]) {
	case 1:
		invalidArgument(Parms, args[*i])
	case 2:
		parseFlag(args, *i, Parms)
	default:
		parseParameter(args, i, Parms)
	}
}

func parseFlag(args []string, i int, Parms *structures.Parms) {
	if args[i][1] == '-' {
		invalidArgument(Parms, args[i])
	}

	switch args[i][1] {
	case 'a':
		Parms.All = true
	case 'c':
		Parms.Checkpoints = true
	case 'd':
		Parms.Debug = true
	case 'f':
		Parms.Force = true
	case 'h':
		Parms.Help = true
	case 'i':
		Parms.Install = true
	case 'l':
		Parms.Local = true
	case 's':
		Parms.Verbose = 0
	case 't':
		Parms.Tolerant = true
	case 'v':
		Parms.Verbose = 3
	case 'V':
		Parms.Verbose = 7
	default:
		invalidArgument(Parms, args[i])
	}
}

func parseParameter(args []string, i *int, Parms *structures.Parms) {
	if args[*i][1] != '-' {
		invalidArgument(Parms, args[*i])
	}
	if *i+1 >= len(args) {
		missingParameterValue(Parms, args[*i])
	}

	name := args[*i][2:]
	(*i)++
	value := args[*i]

	validateParameter(Parms, name, value)
}

func validateParameter(Parms *structures.Parms, name string, value string) {
	switch name {
	case "exclude":
		processExclusions(Parms, value)
	case "format":
		Parms.Format = value
	case "message":
		Parms.Message = value
	case "path":
		Parms.Path = value
	default:
		invalidArgument(Parms, "--"+name)
	}
}

func invalidArgument(Parms *structures.Parms, argument string) {
	cli.Error(RC.Error, *Parms, "Argumento no válido: %q", argument)
}

func missingParameterValue(Parms *structures.Parms, parameter string) {
	cli.Error(RC.Error, *Parms, "Falta el valor del parámetro: %q", parameter)
}

// Prepare discovers the effective Git repositories.
func Prepare(Parms *structures.Parms) {
	processTargets(Parms)
}

func processExclusions(Parms *structures.Parms, values string) {
	for _, value := range strings.Split(values, ",") {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if info, err := os.Stat(value); err == nil && !info.IsDir() {
			addExclusionsFile(Parms, value)
			continue
		}

		Parms.Exclusions = append(Parms.Exclusions, value)
	}

	Parms.Exclusions = uniqueStrings(Parms.Exclusions)
}

func addExclusionsFile(Parms *structures.Parms, path string) {
	file, err := os.Open(path)
	if err != nil {
		cli.Error(RC.Error, *Parms, "No se puede leer el fichero de exclusiones: %q", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		if value == "" {
			continue
		}
		Parms.Exclusions = append(Parms.Exclusions, value)
	}

	if err := scanner.Err(); err != nil {
		cli.Error(RC.Error, *Parms, "Error leyendo el fichero de exclusiones: %q", path)
	}
}

// processTargets resolves roots and discovers Git repositories recursively.
// If a target is inside a repository, that containing repository is selected.
func processTargets(Parms *structures.Parms) {
	roots := Parms.Targets
	if len(roots) == 0 {
		roots = []string{"."}
	}

	Parms.Targets = []string{}
	Parms.Repos = []string{}
	Parms.BlackList = []string{}

	for _, root := range roots {
		path, err := filepath.Abs(root)
		if err != nil {
			cli.Warning(*Parms, "Se ignora %q: no se puede resolver la ruta.", root)
			continue
		}

		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			cli.Warning(*Parms, "Se ignora %q: no existe o no es un directorio.", root)
			continue
		}

		path = filepath.Clean(path)
		Parms.Targets = append(Parms.Targets, path)

		if repository := tools.FindRepo(path); repository != "" {
			Parms.Repos = append(Parms.Repos, repository)
			continue
		}

		discoverRepos(Parms, path)
	}

	Parms.Targets = uniqueStrings(Parms.Targets)
	Parms.Repos = uniqueStrings(Parms.Repos)
}

// discoverRepos recursively discovers Git repositories below path.
func discoverRepos(Parms *structures.Parms, path string) {
	if isExcluded(Parms, filepath.Base(path)) {
		return
	}

	if tools.IsRepo(path) {
		Parms.Repos = append(Parms.Repos, filepath.Clean(path))
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		cli.Warning(*Parms, "Se ignora %q: no se puede leer.", path)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || isExcluded(Parms, entry.Name()) {
			continue
		}
		discoverRepos(Parms, filepath.Join(path, entry.Name()))
	}
}

func isExcluded(Parms *structures.Parms, name string) bool {
	for _, exclusion := range Parms.Exclusions {
		if name == exclusion {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	unique := []string{}
	seen := map[string]bool{}

	for _, value := range values {
		key := filepath.Clean(value)
		if seen[key] {
			continue
		}

		seen[key] = true
		unique = append(unique, value)
	}

	return unique
}
