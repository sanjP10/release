package cmd

import (
	"bitbucket.org/cloudreach/release/bitbucket"
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"net/http"
	"os"
	"strings"
)

// Create for create sub command
type Create struct {
	username string
	password string
	tag      string
	repo     string
	hash     string
}

// Name of sub command
func (*Create) Name() string { return "create" }

// Synopsis of sub command
func (*Create) Synopsis() string { return "create release for bitbucket repo." }

// Usage of sub command
func (*Create) Usage() string {
	return `create [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>]:
  creates tag against bitbucket repo
`
}

// SetFlags required for create sub command
func (c *Create) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.username, "username", "", "username")
	f.StringVar(&c.password, "password", "", "password")
	f.StringVar(&c.repo, "repo", "", "repo")
	f.StringVar(&c.tag, "tag", "", "tag")
	f.StringVar(&c.hash, "hash", "", "hash")
}

// Execute flow for create sub command
func (c *Create) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exit := subcommands.ExitSuccess
	errors := checkCreateFlags(c)
	if len(errors) > 0 {
		exit = subcommands.ExitUsageError
		_, err := os.Stderr.WriteString("missing flags for create:\n" + strings.Join(errors, "\n"))
		if err != nil {
			fmt.Println("Cannot write to stderr", err)
			exit = subcommands.ExitFailure
		}
	} else {
		client := &http.Client{}
		success := bitbucket.CreateTag(c.username, c.password, c.repo, c.tag, c.hash, *client)
		if !success {
			exit = subcommands.ExitFailure
		}
	}
	return exit
}

func checkCreateFlags(c *Create) []string {
	var errors []string
	if len(c.username) == 0 {
		errors = append(errors, "-username required")
	}
	if len(c.password) == 0 {
		errors = append(errors, "-password required")
	}
	if len(c.repo) == 0 {
		errors = append(errors, "-repo required")
	}
	if len(c.tag) == 0 {
		errors = append(errors, "-tag required")
	}
	if len(c.hash) == 0 {
		errors = append(errors, "-hash required")
	}
	return errors
}
