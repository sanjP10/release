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

// Validate for validate sub command
type Validate struct {
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

// Name of subcommand
func (*Validate) Name() string { return "validate" }

// Synopsis of subcommand
func (*Validate) Synopsis() string { return "Validates tag and release for repo to be created." }

// Usage of sub command
func (*Validate) Usage() string {
	return "Validates tag and release for repo to be created.\n"
}

// SetFlags required for create sub command
func (v *Validate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.username, "username", "", "Username (gitlab provider does not require this field). If using ssh provide a username is not git")
	f.StringVar(&v.password, "password", "", "Password or API token (gitlab provider requires an api token). If using a ssh key please provide the password for your ssh key if password protected")
	f.StringVar(&v.email, "email", "", "Required when the provider flag is not supplied, the email for tag")
	f.StringVar(&v.repo, "repo", "", "The repo name, this should include the organisation or owner, required when a provider is supplied")
	f.StringVar(&v.changelog, "changelog", "", "Location of changelog markdown file")
	f.StringVar(&v.hash, "hash", "", "The Full commit hash")
	f.StringVar(&v.host, "host", "", "The host for self hosted instances of the allowed providers")
	f.StringVar(&v.origin, "origin", "", "HTTPS or SSH origin of git repository, to be provided when the provider flag is not provided")
	f.StringVar(&v.provider, "provider", "", "The Git provider, options are github, gitlab or bitbucket, when providing this flag you will be using their APIs")
	f.StringVar(&v.ssh, "ssh", "", "SSH private key file location, please provide Username and password of the SSH file if required. Username defaults to git. This is to be used when the provider flag is not provided")
}

// Execute flow of subcommand
func (v *Validate) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exit := subcommands.ExitSuccess
	errors := checkValidateFlags(v)
	if len(errors) > 0 {
		errors = append(errors, "\n")
		exit = subcommands.ExitUsageError
		_, err := os.Stderr.WriteString("missing flags for validate:\n" + strings.Join(errors, "\n"))
		if err != nil {
			panic("Cannot write to stderr")
		}
	} else {
		changelogFile, err := changelog.ReadChangelogAsString(v.changelog)
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
				desiredTag := changelogObj.ConvertToDesiredTag()
				success, err := validateProviderTag(v, desiredTag, changelogObj)
				if err != nil {
					_, err := os.Stderr.WriteString("Error validating tag with repo " + v.origin + " " + err.Error() + "\n")
					if err != nil {
						panic("Cannot write to stderr")
					}
					exit = subcommands.ExitFailure
				} else if !success {
					_, err := os.Stderr.WriteString("Tag " + strings.TrimSpace(desiredTag) + " already exists\n")
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

func checkValidateFlags(v *Validate) []string {
	var errors []string
	if len(v.provider) == 0 {
		// Use regular git. Check for origin, username/ssh and email
		if len(v.origin) == 0 {
			errors = append(errors, "-origin required")
		}
		if len(v.username) == 0 && len(v.ssh) == 0 {
			errors = append(errors, "-username or -ssh required, for CodeCommit or GCP Source repositories both are required")
		}
		if len(v.password) == 0 && len(v.ssh) == 0 {
			errors = append(errors, "-password required")
		}
		if len(v.email) == 0 {
			errors = append(errors, "-email required")
		}
	} else if ValidProvider(v.provider) {
		// for valid providers check for username and repo
		if len(v.username) == 0 && strings.ToLower(v.provider) != "gitlab" {
			errors = append(errors, "-username required")
		}
		if len(v.password) == 0 {
			errors = append(errors, "-password required")
		}
		if len(v.repo) == 0 {
			errors = append(errors, "-repo required")
		}
	} else {
		// valid provider values
		errors = append(errors, "-provider valid values are "+strings.Join(providers[:], ", "))
	}
	// changelog and hash are mandatory
	if len(v.changelog) == 0 {
		errors = append(errors, "-changelog required")
	}
	if len(v.hash) == 0 {
		errors = append(errors, "-hash required")
	}
	return errors
}

func validateProviderTag(v *Validate, desiredTag string, changelogObj changelog.Properties) (bool, error) {
	success := false
	validTagState := tag.ValidTagState{}
	properties := tag.RepoProperties{
		Password: v.password,
		Tag:      strings.TrimSpace(desiredTag),
		Hash:     v.hash,
		Body:     changelogObj.Changes,
	}
	switch strings.ToLower(v.provider) {
	case "github":
		provider := github.Properties{Username: v.username, Repo: v.repo, Host: v.host, RepoProperties: properties}
		validTagState = provider.ValidateTag()
	case "gitlab":
		provider := gitlab.Properties{Repo: v.repo, Host: v.host, RepoProperties: properties}
		validTagState = provider.ValidateTag()
	case "bitbucket":
		provider := bitbucket.Properties{Username: v.username, Repo: v.repo, Host: v.host, RepoProperties: properties}
		validTagState = provider.ValidateTag()
	default:
		provider := git.Properties{Username: v.username, Email: v.email, Origin: v.origin, SSH: v.ssh, RepoProperties: properties}
		err := provider.InitializeRepository()
		if err != nil {
			return false, err
		}
		validTagState = provider.ValidateTag()
	}
	success = validTagState.TagDoesntExist || validTagState.TagExistsWithProvidedHash
	return success, nil
}
