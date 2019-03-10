package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os"
	"strings"
)

// Validate for validate sub command
type Validate struct {
	username string
	password string
	repo     string
	tag      string
}

// Name of subcommand
func (*Validate) Name() string { return "validate" }

// Synopsis of subcommand
func (*Validate) Synopsis() string { return "validates release version to be created." }

// Usage of subcommand
func (*Validate) Usage() string {
	return `validate [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>]:
  validates tag against bitbucket repo
`
}

// SetFlags for required flags of command
func (v *Validate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.username, "username", "", "username")
	f.StringVar(&v.password, "password", "", "password")
	f.StringVar(&v.repo, "repo", "", "repo")
	f.StringVar(&v.tag, "tag", "", "tag")
}

//Execute flow of subcommand
func (v *Validate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exit := subcommands.ExitSuccess
	errors := checkValidateFlags(v)
	if len(errors) > 0 {
		exit = subcommands.ExitUsageError
		_, err := os.Stderr.WriteString("required flags for validate:\n" + strings.Join(errors, "\n"))
		if err != nil {
			fmt.Println("Cannot write to stderr")
			exit = subcommands.ExitFailure
		}
	}
	return exit
}

func checkValidateFlags(v *Validate) []string {
	var errors []string
	if len(v.username) == 0 {
		errors = append(errors, "-username required")
	}
	if len(v.password) == 0 {
		errors = append(errors, "-password required")
	}
	if len(v.repo) == 0 {
		errors = append(errors, "-repo required")
	}
	if len(v.tag) == 0 {
		errors = append(errors, "-tag required")
	}
	return errors
}
