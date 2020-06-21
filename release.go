package main

import (
	"bitbucket.org/cloudreach/release/internal/commands"
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.Validate{}, "")
	subcommands.Register(&commands.Create{}, "")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
