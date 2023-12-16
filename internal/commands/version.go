package commands

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
	"strings"
)

// Version for validate sub command
type Version struct {
}

// SetFlags of subcommand
func (v *Version) SetFlags(_ *flag.FlagSet) {
}

// Name of subcommand
func (*Version) Name() string { return "version" }

// Synopsis of subcommand
func (*Version) Synopsis() string { return "Version of release tool." }

// Usage of sub command
func (*Version) Usage() string {
	return "Get version of release tool.\n"
}

// Execute flow of subcommand
func (*Version) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exit := subcommands.ExitSuccess
	_, err := os.Stdout.WriteString(strings.TrimSpace("3.2.2") + "\n")
	if err != nil {
		panic("Cannot write to stderr")
	}
	return exit
}
