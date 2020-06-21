package interfaces

// Tagging interface
type Tagging interface {
	create() bool
	validate() bool
}

// RepoProperties properties for repo
type RepoProperties struct {
	Password string
	Repo     string
	Tag      string
	Hash     string
	Host     string
}

// ValidTagState properties for repo
type ValidTagState struct {
	TagDoesntExist            bool
	TagExistsWithProvidedHash bool
}
