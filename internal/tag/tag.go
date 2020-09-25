package tag

// RepoProperties properties for repo
type RepoProperties struct {
	Password string
	Tag      string
	Hash     string
	Body     string
}

// ValidTagState properties for repo
type ValidTagState struct {
	TagDoesntExist            bool
	TagExistsWithProvidedHash bool
}
