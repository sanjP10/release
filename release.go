package main

import (
	"bitbucket.org/cloudreach/release/cmd"
	"bitbucket.org/cloudreach/release/changelog"
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

func main_h() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmd.Validate{}, "")
	subcommands.Register(&cmd.Create{}, "")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

func main() {
	changelogFile, _ := changelog.ReadChangelogAsString("CHANGELOG.md")
	changelogObj := changelog.Properties{}
	changelogObj.GetVersions(changelogFile)
	validateSemantics := changelogObj.ValidateVersionSemantics()
	if validateSemantics {
		changelogObj.RetrieveChanges(changelogFile)
	}
}