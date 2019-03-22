package cmd

import (
	"bitbucket.org/cloudreach/release/bitbucket"
	"context"
	"flag"
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
	hash     string
	host     string
}

// Name of subcommand
func (*Validate) Name() string { return "validate" }

// Synopsis of subcommand
func (*Validate) Synopsis() string { return "validates release version to be created." }

// Usage of subcommand
func (*Validate) Usage() string {
	return `validate [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>] [-host <host> (optional)]:
  validates tag against bitbucket repo
`
}

// SetFlags for required flags of command
func (v *Validate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.username, "username", "", "username")
	f.StringVar(&v.password, "password", "", "password")
	f.StringVar(&v.repo, "repo", "", "repo")
	f.StringVar(&v.tag, "tag", "", "tag")
	f.StringVar(&v.hash, "hash", "", "hash")
}

//Execute flow of subcommand
func (v *Validate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exit := subcommands.ExitSuccess
	errors := checkValidateFlags(v)
	if len(errors) > 0 {
		exit = subcommands.ExitUsageError
		_, err := os.Stderr.WriteString("missing flags for validate:\n" + strings.Join(errors, "\n"))
		if err != nil {
			panic("Cannot write to stderr")
		}
	} else {
		tag := bitbucket.RepoProperties{Username: v.username, Password: v.password, Tag: v.tag, Repo: v.repo, Hash: v.hash, Host: v.host}
		success := tag.ValidateTag()
		if !success {
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
	if len(v.hash) == 0 {
		errors = append(errors, "-hash required")
	}
	return errors
}
