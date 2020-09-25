package git

import (
	"bitbucket.org/cloudreach/release/internal/tag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
	"time"
)

// Properties for git repo
type Properties struct {
	tag.RepoProperties
	Username string
	Email    string
	Origin   string
	SSH      string
}

var repository *git.Repository

// InitializeRepository initializes the repository sets up origins and fetches
func (r *Properties) InitializeRepository() error {
	var err error
	repository, err = git.Init(memory.NewStorage(), nil)
	if err != nil {
		fmt.Println("Error Initializing repository", err)
		return err
	}
	_, err = repository.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{r.Origin},
	})
	if err != nil {
		fmt.Println("Error Setting origin for repository", err)
		return err
	}
	auth, err := getAuth(r.SSH, r.Username, r.Password)
	if err != nil {
		return err
	}
	err = repository.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"+refs/tags/*:refs/tags/*", "+refs/heads/*:refs/remotes/origin/*"},
		Auth:       auth,
	})

	if err != nil {
		fmt.Println("Error Fetching repository", err)
	}
	return err
}

//ValidateTag checks a tag does not exist or has the same hash
func (r *Properties) ValidateTag() tag.ValidTagState {
	validTag := tag.ValidTagState{TagDoesntExist: false, TagExistsWithProvidedHash: false}
	tagref, err := repository.Tag(r.Tag)
	if err != nil {
		if err.Error() == "tag not found" {
			validTag.TagDoesntExist = true
		}
	} else {
		tagObject, err := repository.TagObject(tagref.Hash())
		if err != nil {
			fmt.Println("Error retrieving tag details", err)
		}
		if tagObject.Target.String() == r.Hash {
			validTag.TagExistsWithProvidedHash = true
		}
	}
	return validTag
}

// CreateTag creates a git tag
func (r *Properties) CreateTag() bool {
	createTag := false
	validTagState := r.ValidateTag()
	if validTagState.TagExistsWithProvidedHash {
		createTag = true
	} else if validTagState.TagDoesntExist {
		_, err := repository.CreateTag(r.Tag, plumbing.NewHash(r.Hash), &git.CreateTagOptions{
			Tagger: &object.Signature{
				Name:  r.Username,
				Email: r.Email,
				When:  time.Time{},
			},
			Message: r.Body,
		})
		if err != nil {
			fmt.Println("Error Creating tag", err)
			return createTag
		}
		auth, err := getAuth(r.SSH, r.Username, r.Password)
		if err != nil {
			return createTag
		}
		po := &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/" + r.Tag + ":refs/tags/" + r.Tag)},
			Auth:       auth,
		}
		err = repository.Push(po)
		if err != nil {
			fmt.Println("Error Pushing tag", err)
			return createTag
		}
		createTag = true
	}
	return createTag
}

func getAuth(filePath string, username string, password string) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	var err error
	if len(filePath) > 0 {
		if username == "" {
			username = "git"
		}
		auth, err = ssh.NewPublicKeysFromFile(username, filePath, password)
		if err != nil {
			fmt.Println("Error Setting SSH Key", err)
			return nil, err
		}
	} else {
		auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}
	return auth, err
}
