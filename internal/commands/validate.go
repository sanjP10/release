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
	f.StringVar(&v.username, "username", "", "username (gitlab does not require this field)")
	f.StringVar(&v.password, "password", "", "password or api token (gitlab requires an api token)")
	f.StringVar(&v.email, "email", "", "Required when a provider is not supplied, the email for tag")
	f.StringVar(&v.repo, "repo", "", "repo name, required when a provider is supplided")
	f.StringVar(&v.changelog, "changelog", "", "location of changelog markdown file")
	f.StringVar(&v.hash, "hash", "", "full commit hash")
	f.StringVar(&v.host, "host", "", "host override for provider specific APIs")
	f.StringVar(&v.origin, "origin", "", "origin of git repository")
	f.StringVar(&v.provider, "provider", "", "git provider, options are github, gitlab or bitbucket")
}

//Execute flow of subcommand
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
	if len(v.username) == 0 && v.provider != "gitlab" {
		errors = append(errors, "-username required")
	}
	if len(v.password) == 0 {
		errors = append(errors, "-password required")
	}
	if len(v.provider) > 0 && len(v.repo) == 0 {
		errors = append(errors, "-repo required")
	}
	if len(v.changelog) == 0 {
		errors = append(errors, "-changelog required")
	}
	if len(v.hash) == 0 {
		errors = append(errors, "-hash required")
	}
	if len(v.provider) > 0 && !ValidProvider(v.provider) {
		errors = append(errors, "-provider required, valid values are "+strings.Join(providers[:], ", "))
	}
	if len(v.provider) == 0 {
		if len(v.email) == 0 {
			errors = append(errors, "-email required")
		}
		if len(v.origin) == 0 {
			errors = append(errors, "-origin required")
		}
	}
	return errors
}

func validateProviderTag(v *Validate, desiredTag string, changelogObj changelog.Properties) (bool, error) {
	success := false
	validTagState := tag.ValidTagState{}
	properties := tag.RepoProperties{
		Password: v.password,
		Tag:      strings.TrimSpace(desiredTag),
		Hash:     v.hash}
	switch strings.ToLower(v.provider) {
	case "github":
		provider := github.Properties{Username: v.username, Body: changelogObj.Changes, Repo: v.repo, Host: v.host, RepoProperties: properties}
		validTagState = provider.ValidateTag()
	case "gitlab":
		provider := gitlab.Properties{Body: changelogObj.Changes, Repo: v.repo, Host: v.host, RepoProperties: properties}
		validTagState = provider.ValidateTag()
	case "bitbucket":
		provider := bitbucket.Properties{Username: v.username, Repo: v.repo, Host: v.host, RepoProperties: properties}
		validTagState = provider.ValidateTag()
	default:
		provider := git.Properties{Username: v.username, Email: v.email, Body: changelogObj.Changes, Origin: v.origin, RepoProperties: properties}
		err := provider.InitializeRepository()
		if err != nil {
			return false, err
		}
		validTagState = provider.ValidateTag()
	}
	success = validTagState.TagDoesntExist || validTagState.TagExistsWithProvidedHash
	return success, nil
}
