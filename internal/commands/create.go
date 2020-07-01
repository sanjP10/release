package commands

import (
	"bitbucket.org/cloudreach/release/internal/changelog"
	"bitbucket.org/cloudreach/release/internal/tag"
	"bitbucket.org/cloudreach/release/internal/tag/providers/bitbucket"
	"bitbucket.org/cloudreach/release/internal/tag/providers/git"
	"bitbucket.org/cloudreach/release/internal/tag/providers/github"
	"bitbucket.org/cloudreach/release/internal/tag/providers/gitlab"
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
	email     string
	changelog string
	repo      string
	hash      string
	host      string
	origin    string
	provider  string
}

// Name of sub command
func (*Create) Name() string { return "create" }

// Synopsis of sub command
func (*Create) Synopsis() string { return "Creates tag and release for repo." }

// Usage of sub command
func (*Create) Usage() string {
	return "Creates tag and release for repo.\n"
}

// SetFlags required for create sub command
func (c *Create) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.username, "username", "", "username (gitlab does not require this field)")
	f.StringVar(&c.password, "password", "", "password or api token (gitlab requires an api token)")
	f.StringVar(&c.password, "email", "", "Required when a provider is not supplied, the email for tag")
	f.StringVar(&c.repo, "repo", "", "repo name, required when a provider is supplided")
	f.StringVar(&c.changelog, "changelog", "", "location of changelog markdown file")
	f.StringVar(&c.hash, "hash", "", "full commit hash")
	f.StringVar(&c.host, "host", "", "host override for provider specific APIs")
	f.StringVar(&c.host, "origin", "", "origin of git repository")
	f.StringVar(&c.provider, "provider", "", "git provider, options are github, gitlab or bitbucket")
}

// Execute flow for create sub command
func (c *Create) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exit := subcommands.ExitSuccess
	errors := checkCreateFlags(c)
	if len(errors) > 0 {
		errors = append(errors, "\n")
		exit = subcommands.ExitUsageError
		_, err := os.Stderr.WriteString("missing flags for create:\n" + strings.Join(errors, "\n"))
		if err != nil {
			panic("Cannot write to stderr")
		}
	} else {
		changelogFile, err := changelog.ReadChangelogAsString(c.changelog)
		if err != nil {
			exit = subcommands.ExitUsageError
			_, err := os.Stderr.WriteString("Unable to read changelog\n")
			if err != nil {
				panic("Cannot write to stderr")
			}
		} else {
			changelogObj := changelog.Properties{}
			changelogObj.GetVersions(changelogFile)
			validSemantics := changelogObj.ValidateVersionSemantics()
			if !validSemantics {
				exit = subcommands.ExitFailure
				_, err := os.Stderr.WriteString("Invalid version semantics\n")
				if err != nil {
					panic("Cannot write to stderr")
				}
			} else {
				changelogObj.RetrieveChanges(changelogFile)
				desiredTag := changelogObj.ConvertToDesiredTag()
				success := createProviderTag(c, desiredTag, changelogObj)
				if !success {
					_, err := os.Stderr.WriteString("Error creating Tag " + strings.TrimSpace(desiredTag) + "\n")
					if err != nil {
						panic("Cannot write to stderr")
					}
					exit = subcommands.ExitFailure
				} else {
					_, err := os.Stdout.WriteString(strings.TrimSpace(desiredTag) + "\n")
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
	if len(c.username) == 0 && c.provider != "gitlab" {
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
	if len(c.provider) > 0 && !ValidProvider(c.provider) {
		errors = append(errors, "-provider required, valid values are "+strings.Join(providers[:], ", "))
	}
	if len(c.provider) == 0 {
		if len(c.email) == 0 {
			errors = append(errors, "-email required")
		}
		if len(c.origin) == 0 {
			errors = append(errors, "-origin required")
		}
	}
	return errors
}

func createProviderTag(c *Create, desiredTag string, changelogObj changelog.Properties) bool {
	success := false
	properties := tag.RepoProperties{
		Password: c.password,
		Tag:      strings.TrimSpace(desiredTag),
		Hash:     c.hash}
	switch strings.ToLower(c.provider) {
	case "github":
		provider := github.Properties{Username: c.username, Body: changelogObj.Changes, Repo: c.repo, Host: c.host, RepoProperties: properties}
		success = provider.CreateTag()
	case "gitlab":
		provider := gitlab.Properties{Body: changelogObj.Changes, Repo: c.repo, Host: c.host, RepoProperties: properties}
		success = provider.CreateTag()
	case "bitbucket":
		provider := bitbucket.Properties{Username: c.username, Repo: c.repo, Host: c.host, RepoProperties: properties}
		success = provider.CreateTag()
	default:
		provider := git.Properties{Username: c.username, Email: c.email, Body: changelogObj.Changes, Origin: c.origin, RepoProperties: properties}
		err := provider.InitializeRepository()
		if err != nil {
			_, err := os.Stderr.WriteString("Error initializing repository " + c.origin + "\n")
			if err != nil {
				panic("Cannot write to stderr")
			}
			return false
		}
		success = provider.CreateTag()
	}
	return success
}
