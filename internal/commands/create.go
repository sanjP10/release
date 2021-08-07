package commands

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/sanjP10/release/internal/changelog"
	"github.com/sanjP10/release/internal/tag"
	"github.com/sanjP10/release/internal/tag/providers/bitbucket"
	"github.com/sanjP10/release/internal/tag/providers/git"
	"github.com/sanjP10/release/internal/tag/providers/github"
	"github.com/sanjP10/release/internal/tag/providers/gitlab"
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
	ssh       string
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
	f.StringVar(&c.username, "username", "", "Username (gitlab provider does not require this field). If using ssh provide a username is not git")
	f.StringVar(&c.password, "password", "", "Password or API token (gitlab provider requires an api token). If using a ssh key please provide the password for your ssh key if password protected")
	f.StringVar(&c.email, "email", "", "Required when the provider flag is not supplied, the email for tag")
	f.StringVar(&c.repo, "repo", "", "The repo name, this should include the organisation or owner, required when a provider is supplied")
	f.StringVar(&c.changelog, "changelog", "", "Location of changelog markdown file")
	f.StringVar(&c.hash, "hash", "", "The Full commit hash")
	f.StringVar(&c.host, "host", "", "The host for self hosted instances of the allowed providers")
	f.StringVar(&c.origin, "origin", "", "Https or SSH origin of git repository, to be provided when the provider flag is not provided")
	f.StringVar(&c.provider, "provider", "", "The Git provider, options are github, gitlab or bitbucket, when providing this flag you will be using their API's")
	f.StringVar(&c.ssh, "ssh", "", "SSH private key file location, please provide Username and password of the SSH file if required. Username defaults to git. This is to be used when the provider flag is not provided")
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
				success, err := createProviderTag(c, desiredTag, changelogObj)
				if err != nil {
					_, err := os.Stderr.WriteString("Error creating tag with repo " + c.origin + " " + err.Error() + "\n")
					if err != nil {
						panic("Cannot write to stderr")
					}
					exit = subcommands.ExitFailure
				} else if !success {
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
	if len(c.provider) == 0 {
		// Use regular git. Check for origin, username/ssh and email
		if len(c.origin) == 0 {
			errors = append(errors, "-origin required")
		}
		if len(c.username) == 0 && len(c.ssh) == 0 {
			errors = append(errors, "-username or -ssh required, for CodeCommit or GCP Source repositories both are required")
		}
		if len(c.password) == 0 && len(c.ssh) == 0 {
			errors = append(errors, "-password required")
		}
		if len(c.email) == 0 {
			errors = append(errors, "-email required")
		}
	} else if ValidProvider(c.provider) {
		// for valid providers check for username and repo
		if len(c.username) == 0 && strings.ToLower(c.provider) != "gitlab" {
			errors = append(errors, "-username required")
		}
		if len(c.password) == 0 {
			errors = append(errors, "-password required")
		}
		if len(c.repo) == 0 {
			errors = append(errors, "-repo required")
		}
	} else {
		// valid provider values
		errors = append(errors, "-provider valid values are "+strings.Join(providers[:], ", "))
	}
	// changelog and hash are mandatory
	if len(c.changelog) == 0 {
		errors = append(errors, "-changelog required")
	}
	if len(c.hash) == 0 {
		errors = append(errors, "-hash required")
	}
	return errors
}

func createProviderTag(c *Create, desiredTag string, changelogObj changelog.Properties) (bool, error) {
	success := false
	properties := tag.RepoProperties{
		Password: c.password,
		Tag:      strings.TrimSpace(desiredTag),
		Hash:     c.hash,
		Body:     changelogObj.Changes,
	}
	switch strings.ToLower(c.provider) {
	case "github":
		provider := github.Properties{Username: c.username, Repo: c.repo, Host: c.host, RepoProperties: properties}
		success = provider.CreateTag()
	case "gitlab":
		provider := gitlab.Properties{Repo: c.repo, Host: c.host, RepoProperties: properties}
		success = provider.CreateTag()
	case "bitbucket":
		provider := bitbucket.Properties{Username: c.username, Repo: c.repo, Host: c.host, RepoProperties: properties}
		success = provider.CreateTag()
	default:
		provider := git.Properties{Username: c.username, Email: c.email, Origin: c.origin, SSH: c.ssh, RepoProperties: properties}
		err := provider.InitializeRepository()
		if err != nil {
			return false, err
		}
		success = provider.CreateTag()
	}
	return success, nil
}
