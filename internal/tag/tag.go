package tag

// Tagging interface
type Tagging interface {
	create() bool
	validate() bool
}

// RepoProperties properties for repo
type RepoProperties struct {
	Password string
	Tag      string
	Hash     string
}

// ValidTagState properties for repo
type ValidTagState struct {
	TagDoesntExist            bool
	TagExistsWithProvidedHash bool
}
