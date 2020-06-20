package tagging

// Tagging interface
type Tagging interface {
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

// GithubProperties properties for repo
type GithubProperties struct {
	RepoProperties
}

// GitlabProperties properties for repo
type GitlabProperties struct {
	RepoProperties
}

// BitbucketProperties properties for repo
type BitbucketProperties struct {
	RepoProperties
}

// ValidTagState properties for repo
type ValidTagState struct {
	TagDoesntExist            bool
	TagExistsWithProvidedHash bool
}
