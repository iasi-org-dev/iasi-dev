package main

import "iasi-dev/internal/cli"

func printHelp() {
	cli.Direct(`IASI Dev
	
	Usage:
	  iasi-dev <command> [-h] [-s] [-v|-V] [-a] [-c] [-d] [-f] [-l] [-t] [-i] [--path value] [--exclude value[,value]*] [--format value] [--message value] [target...]
	  iasi-dev workflow <build|publish|release> [-h] [-s] [-v|-V] [-a] [-c] [-d] [-f] [-l] [-t] [-i] [--path value] [--exclude value[,value]*] [--format value] [--message value] [target...]
	  iasi-dev workflow promote vMAJOR.MINOR.PATCH [-h] [-s] [-v|-V] [-d] [--path value] [target...]
	  iasi-dev promote vMAJOR.MINOR.PATCH [-h] [-s] [-v|-V] [-d] [--path value] [target...]
	  iasi-dev restore [vMAJOR.MINOR.PATCH] [-h] [-s] [-v|-V] [-d] [--path value] [target...]
	  iasi-dev version [organization] [-h] [-s] [-v|-V] [-d] [--path value]
	
	Commands:
	  help       Show help
	  build      Build through iasi.quarto
	  publish    Publish through iasi.quarto
	  release    Release through iasi.quarto
	  commit     Commit
	  promote    Validate and create a new stable organization version (promotion steps pending)
	  restore    Restore repositories to a tagged version; without a version, restore main
	  workflow   Run a workflow repository by repository
	  sync       Sync shared files from iasi-common
	  version    Show VERSION for the explicit or current organization
	
	Workflows:
	  build      Build and commit
	  publish    Publish and commit; optionally include previous stages with -a
	  release    Release and commit; optionally include previous stages with -a
	  promote    Promote and publish the stable organization to GitHub (implementation pending)
	
	Options:
	  -h         Show help.
	  -s         Silent output.
	  -v         Verbose output.
	  -V         Very verbose output; preserves process RC 1 (NothingToDo).
	  -a         Execute previous workflow stages too.
	  -c         Use intermediate workflows as checkpoints.
	  -d         Show debug messages.
	  -f         Force the operation when supported.
	  -i         Install the artifact when applicable.
	  -l         Commit locally without push.
	  -t         Continue when an operation fails.
	
	Parameters:
	  --path value               Change to this directory before common preparation.
	  --exclude value[,value]*  Add exclusions. Existing files are read one exclusion per line; .git, .github and tests are always excluded.
	  --format value            Output format passed to build.
	  --message value           Commit message.
	`)
}
