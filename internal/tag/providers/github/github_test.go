package github

import (
	"github.com/sanjP10/release/internal/tag"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestValidateTagNotExisting(t *testing.T) {
	// Testing tag not existing
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusNotFound)
	assertTest := assert.New(t)
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash"}}
	results := repo.ValidateTag()
	assertTest.True(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingSameHash(t *testing.T) {
	target := Object{Sha: "hash"}
	response := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusOK).
		JSON(response)

	assertTest := assert.New(t)
	// Testing 200 response and hash is the same
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash"}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.True(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingMismatchHash(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response but hash is not the same
	target := Object{Sha: "hash"}
	response := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusOK).
		JSON(response)
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "not_hash"}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagOtherError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing a 403
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash"}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestCreateTagNotFound(t *testing.T) {
	// Testing tag not existing
	target := Object{Sha: "tag"}
	response := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusNotFound).
		JSON(response)

	assertTest := assert.New(t)
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash"}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagSuccessful(t *testing.T) {
	// Testing 201 response
	body := Release{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusCreated).
		JSON(body)
	assertTest := assert.New(t)
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash", Body: "hello"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagSuccessfulWithHostOverride(t *testing.T) {
	// Testing 201 response
	body := Release{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.personal-github.com").
		Get("/api/v3/repos/repo/git/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.personal-github.com").
		Post("/api/v3/repos/repo/releases").
		Reply(http.StatusCreated).
		JSON(body)
	assertTest := assert.New(t)
	repo := Properties{Username: "username", Repo: "repo", Host: "https://api.personal-github.com", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash", Body: "hello"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagAlreadyExists(t *testing.T) {
	target := Object{Sha: "hash"}
	response := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK).
		JSON(response)
	assertTest := assert.New(t)
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "test", Hash: "hash", Body: "hello"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateError(t *testing.T) {
	target := Object{Sha: "hash"}
	response := Tag{Object: target}
	errorMessage := Error{Code: "already_exists"}
	errorResponse := BadResponse{Errors: []Error{errorMessage}}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusNotFound).
		JSON(response)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusUnprocessableEntity).
		JSON(errorResponse)
	assertTest := assert.New(t)
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{"password", "test", "hash", "hello"}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherError(t *testing.T) {
	errorMessage := Error{Code: "blah"}
	response := BadResponse{Errors: []Error{errorMessage}}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusOK)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusBadRequest).
		JSON(response)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash", Body: "hello"}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherErrorResponse(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags/tag").
		Reply(http.StatusOK)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := Properties{Username: "username", Repo: "repo", Host: "", RepoProperties: tag.RepoProperties{Password: "password", Tag: "tag", Hash: "hash", Body: "hello"}}
	assertTest.False(repo.CreateTag())
}
