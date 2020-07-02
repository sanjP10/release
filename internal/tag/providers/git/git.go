package git

import (
	"bitbucket.org/cloudreach/release/internal/tag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
	"time"
)

type Properties struct {
	tag.RepoProperties
	Username string
	Email    string
	Body     string
	Origin   string
}

var repository *git.Repository

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
	err = repository.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec("+refs/tags/*:refs/tags/*")},
		Auth: &http.BasicAuth{
			Username: r.Username,
			Password: r.Password,
		},
	})

	if err != nil {
		fmt.Println("Error Fetching repository", err)
	}
	return err
}

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
		po := &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/" + r.Tag + ":refs/tags/" + r.Tag)},
			Auth: &http.BasicAuth{
				Username: r.Username,
				Password: r.Password,
			},
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
