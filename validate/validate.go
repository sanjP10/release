package validate

import (
	"context"
	"flag"
	"github.com/google/subcommands"
)

type Validate struct {
	username string
	password string
	repo string
	tag string
}

func (*Validate) Name() string     { return "validate" }
func (*Validate) Synopsis() string { return "validates release version to be created." }
func (*Validate) Usage() string {
	return `validate [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>]:
  validates tag against bitbucket repo
`
}

func (v *Validate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.username, "username","", "username")
	f.StringVar(&v.password, "password", "", "password")
	f.StringVar(&v.repo, "repo", "","repo")
	f.StringVar(&v.tag, "tag", "", "tag")
}


func (v *Validate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
