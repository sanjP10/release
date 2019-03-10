package main

import (
	"bitbucket.org/cloudreach/release/cmd"
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmd.Validate{}, "")
	subcommands.Register(&cmd.Create{}, "")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
