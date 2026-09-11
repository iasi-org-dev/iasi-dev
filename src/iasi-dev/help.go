package main

import "iasi-dev/internal/cli"

func printHelp() {
	cli.Direct(`IASI Dev
	
	Usage:
	  iasi-dev <command> [-h] [-s] [-v|-V] [-a] [-c] [-d] [-f] [-l] [-t] [-i] [--exclude value[,value]*] [--format value] [--message value] [target...]
	  iasi-dev workflow <build|publish|release> [-h] [-s] [-v|-V] [-a] [-c] [-d] [-f] [-l] [-t] [-i] [--exclude value[,value]*] [--format value] [--message value] [target...]
	
	Commands:
	  help       Show help
	  build      Build through iasi.quarto
	  publish    Publish through iasi.quarto
	  release    Release through iasi.quarto
	  commit     Commit
	  workflow   Run a workflow repository by repository
	  sync       Sync shared files from iasi-common
	
	Options:
	  -h         Show help.
	  -s         Silent output.
	  -v         Verbose output.
	  -V         Very verbose output.
	  -a         Execute previous workflow stages too.
	  -c         Use intermediate workflows as checkpoints.
	  -d         Show debug messages.
	  -f         Force the operation when supported.
	  -i         Install the artifact when applicable.
	  -l         Commit locally without push.
	  -t         Continue when an operation fails.
	
	Parameters:
	  --exclude value[,value]*  Add exclusions. Existing files are read one exclusion per line; .git and tests are always excluded.
	  --format value            Output format passed to build.
	  --message value           Commit message.
	`)
}
