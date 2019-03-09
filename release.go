package main

import (
	"bitbucket.org/cloudreach/release/create"
	"bitbucket.org/cloudreach/release/validate"
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&validate.Validate{}, "")
	subcommands.Register(&create.Create{}, "")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
