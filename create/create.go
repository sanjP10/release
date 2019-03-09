package create

import (
	"context"
	"flag"
	"github.com/google/subcommands"
)

type Create struct {
	username string
	password string
	tag string
	repo string
}

func (*Create) Name() string     { return "create" }
func (*Create) Synopsis() string { return "create release for bitbucket repo." }
func (*Create) Usage() string {
	return `create [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>]:
  creates tag against bitbucket repo
`
}

func (c *Create) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.username, "username","", "username")
	f.StringVar(&c.password, "password", "", "password")
	f.StringVar(&c.repo, "repo", "","repo")
	f.StringVar(&c.tag, "tag", "", "tag")
}


func (c *Create) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
