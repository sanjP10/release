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

// Validate for validate sub command
type Validate struct {
	username  string
	password  string
	repo      string
	changelog string
	hash      string
	host      string
	provider  string
}

// Name of subcommand
func (*Validate) Name() string { return "validate" }

// Synopsis of subcommand
func (*Validate) Synopsis() string { return "validates release version to be created." }

// Usage of subcommand
func (*Validate) Usage() string {
	return `validate [-username <username>] [-password <password/token>] [-repo <repo>] [-changelog <changelog md file>] [-provider <github/gitlab/bitbucket>] [-host <host> (optional)]:
  validates tag against bitbucket repo
`
}

// SetFlags for required flags of command
func (v *Validate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.username, "username", "", "username")
	f.StringVar(&v.password, "password", "", "password")
	f.StringVar(&v.repo, "repo", "", "repo")
	f.StringVar(&v.changelog, "changelog", "", "changelog")
	f.StringVar(&v.hash, "hash", "", "hash")
	f.StringVar(&v.host, "host", "", "host")
	f.StringVar(&v.provider, "provider", "", "provider")
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
		changelogFile, err := changelog.ReadChangelogAsString(v.changelog)
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
				desiredTag := changelogObj.ConvertToDesiredTag()
				success := validateProviderTag(v, desiredTag, changelogObj)
				if !success {
					_, err := os.Stderr.WriteString("Tag cannot be created or already exists")
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
	if len(v.changelog) == 0 {
		errors = append(errors, "-changelog required")
	}
	if len(v.hash) == 0 {
		errors = append(errors, "-hash required")
	}

	if !ValidProvider(v.provider) {
		errors = append(errors, "-provider required, valid values are "+strings.Join(providers[:], ", "))
	}
	return errors
}

func validateProviderTag(v *Validate, desiredTag string, changelogObj changelog.Properties) bool {
	success := false
	if strings.ToLower(v.provider) == "github" {
		tag := github.RepoProperties{
			Username: v.username,
			Password: v.password,
			Repo:     v.repo,
			Tag:      strings.TrimSpace(desiredTag),
			Body:     changelogObj.Changes,
			Hash:     v.hash,
			Host:     v.host}
		success = tag.ValidateTag()
	} else if strings.ToLower(v.provider) == "gitlab" {
		tag := gitlab.RepoProperties{
			Token: v.password,
			Repo:  v.repo,
			Tag:   strings.TrimSpace(desiredTag),
			Body:  changelogObj.Changes,
			Hash:  v.hash,
			Host:  v.host}
		success = tag.ValidateTag()
	} else if strings.ToLower(v.provider) == "bitbucket" {
		tag := bitbucket.RepoProperties{
			Username: v.username,
			Password: v.password,
			Repo:     v.repo,
			Tag:      strings.TrimSpace(desiredTag),
			Hash:     v.hash,
			Host:     v.host}
		success = tag.ValidateTag()
	}
	return success
}
