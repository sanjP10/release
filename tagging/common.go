package tagging

// Bitbucket interface for bitbucket
type Git interface {
	create() bool
	validate() bool
}

// RepoProperties properties for repo
type RepoProperties struct {
	Username string
	Password string
	Repo     string
	Tag      string
	Hash     string
	Host     string
	Body     string
}

type GithubProperties struct {
	RepoProperties
}

type GitlabProperties struct {
	RepoProperties
}

type BitbucketProperties struct {
	RepoProperties
}
