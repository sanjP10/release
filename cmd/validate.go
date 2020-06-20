package cmd

import (
	"bitbucket.org/cloudreach/release/changelog"
	"bitbucket.org/cloudreach/release/tagging"
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
func (*Validate) Synopsis() string { return "Validates tag and release for repo to be created." }

// Usage of sub command
func (*Validate) Usage() string {
	return "Validates tag and release for repo to be created.\n"
}

// SetFlags required for create sub command
func (v *Validate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.username, "username", "", "username (gitlab does not require this field)")
	f.StringVar(&v.password, "password", "", "password or api token (gitlab requires an api token)")
	f.StringVar(&v.repo, "repo", "", "repo name")
	f.StringVar(&v.changelog, "changelog", "", "location of changelog markdown file")
	f.StringVar(&v.hash, "hash", "", "full commit hash")
	f.StringVar(&v.host, "host", "", "host override")
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
				success := validateProviderTag(v, desiredTag, changelogObj)
				if !success {
					_, err := os.Stderr.WriteString("Tag cannot be created or already exists\n")
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
	validTagState := tagging.ValidTagState{}
	properties := tagging.RepoProperties{
		Username: v.username,
		Password: v.password,
		Repo:     v.repo,
		Tag:      strings.TrimSpace(desiredTag),
		Body:     changelogObj.Changes,
		Hash:     v.hash,
		Host:     v.host}
	switch strings.ToLower(v.provider) {
	case "github":
		provider := tagging.GithubProperties{RepoProperties: properties}
		validTagState = provider.ValidateTag()
	case "gitlab":
		provider := tagging.GitlabProperties{RepoProperties: properties}
		validTagState = provider.ValidateTag()
	case "bitbucket":
		provider := tagging.BitbucketProperties{RepoProperties: properties}
		validTagState = provider.ValidateTag()
	}
	success = validTagState.TagDoesntExist || validTagState.TagExistsWithProvidedHash
	return success
}
