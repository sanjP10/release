package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/sanjP10/release/internal/commands"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.Validate{}, "")
	subcommands.Register(&commands.Create{}, "")
	subcommands.Register(&commands.Version{}, "")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
