package cmd

import (
	"bitbucket.org/cloudreach/release/bitbucket"
	"bitbucket.org/cloudreach/release/changelog"
	"bitbucket.org/cloudreach/release/github"
	"bitbucket.org/cloudreach/release/gitlab"
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
	"strings"
)

// Create for create sub command
type Create struct {
	username  string
	password  string
	changelog string
	repo      string
	hash      string
	host      string
	provider  string
}

// Name of sub command
func (*Create) Name() string { return "create" }

// Synopsis of sub command
func (*Create) Synopsis() string { return "create release for bitbucket repo." }

// Usage of sub command
func (*Create) Usage() string {
	return `create [-username <username>] [-password <password/token>] [-repo <repo>] [-changelog <changelog md file>] [-provider <github/gitlab/bitbucket>] [-host <host> (optional)]:
  creates tag against bitbucket repo
`
}

// SetFlags required for create sub command
func (c *Create) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.username, "username", "", "username")
	f.StringVar(&c.password, "password", "", "password")
	f.StringVar(&c.repo, "repo", "", "repo")
	f.StringVar(&c.changelog, "changelog", "", "changelog")
	f.StringVar(&c.hash, "hash", "", "hash")
	f.StringVar(&c.host, "host", "", "host")
	f.StringVar(&c.provider, "provider", "", "provider")
}

// Execute flow for create sub command
func (c *Create) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exit := subcommands.ExitSuccess
	errors := checkCreateFlags(c)
	if len(errors) > 0 {
		exit = subcommands.ExitUsageError
		_, err := os.Stderr.WriteString("missing flags for create:\n" + strings.Join(errors, "\n"))
		if err != nil {
			panic("Cannot write to stderr")
		}
	} else {
		changelogFile, err := changelog.ReadChangelogAsString(c.changelog)
		if err != nil {
			exit = subcommands.ExitUsageError
			_, err := os.Stderr.WriteString("Unable to read changelog")
			if err != nil {
				panic("Cannot write to stderr")
			}
		} else {
			changelogObj := changelog.Properties{}
			changelogObj.GetVersions(changelogFile)
			validSemantics := changelogObj.ValidateVersionSemantics()
			if !validSemantics {
				exit = subcommands.ExitFailure
				_, err := os.Stderr.WriteString("Invalid version semantics")
				if err != nil {
					panic("Cannot write to stderr")
				}
			} else {
				changelogObj.RetrieveChanges(changelogFile)
				desiredTag := changelogObj.ConvertToDesiredTag()
				success := createProviderTag(c, desiredTag, changelogObj)
				if !success {
					_, err := os.Stderr.WriteString("Error creating Tag" + desiredTag)
					if err != nil {
						panic("Cannot write to stderr")
					}
					exit = subcommands.ExitFailure
				} else {
					_, err := os.Stdout.WriteString(strings.TrimSpace(desiredTag))
					if err != nil {
						panic("Cannot write to stderr")
					}
				}
			}
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
	if len(c.changelog) == 0 {
		errors = append(errors, "-changelog required")
	}
	if len(c.hash) == 0 {
		errors = append(errors, "-hash required")
	}

	if !ValidProvider(c.provider) {
		errors = append(errors, "-provider required, valid values are "+strings.Join(providers[:], ", "))
	}
	return errors
}

func createProviderTag(c *Create, desiredTag string, changelogObj changelog.Properties) bool {
	success := false
	provider := strings.ToLower(c.provider)
	if provider == "github" {
		tag := github.RepoProperties{
			Username: c.username,
			Password: c.password,
			Repo:     c.repo,
			Tag:      strings.TrimSpace(desiredTag),
			Body:     changelogObj.Changes,
			Hash:     c.hash,
			Host:     c.host}
		success = tag.CreateTag()
	} else if provider == "gitlab" {
		tag := gitlab.RepoProperties{
			Token: c.password,
			Repo:  c.repo,
			Tag:   strings.TrimSpace(desiredTag),
			Body:  changelogObj.Changes,
			Hash:  c.hash,
			Host:  c.host}
		success = tag.CreateTag()
	} else if provider == "bitbucket" {
		tag := bitbucket.RepoProperties{
			Username: c.username,
			Password: c.password,
			Repo:     c.repo,
			Tag:      strings.TrimSpace(desiredTag),
			Hash:     c.hash,
			Host:     c.host}
		success = tag.CreateTag()
	}
	return success
}
