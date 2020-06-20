package github

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestValidateTagNotExisting(t *testing.T) {
	// Testing tag not existing
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}
	assertTest.True(repo.ValidateTag())
}

func TestValidateTagUnauthorized(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusUnauthorized)
	assertTest := assert.New(t)
	// Testing a 403
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}
	assertTest.False(repo.ValidateTag())
}

func TestValidateTagExistingSameHash(t *testing.T) {
	target := Object{Sha: "hash"}
	tag := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK).
		JSON(tag)

	assertTest := assert.New(t)
	// Testing 200 response and hash is the same
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}
	assertTest.True(repo.ValidateTag())
}

func TestValidateTagExistingMismatchHash(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response but hash is not the same
	target := Object{Sha: "hash"}
	tag := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK).
		JSON(tag)
	repo := RepoProperties{"username", "password", "repo", "tag", "not_hash", "", ""}
	assertTest.False(repo.ValidateTag())
}

func TestValidateTagOtherError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing a 403
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}
	assertTest.False(repo.ValidateTag())
}

func TestCreateTagNotFound(t *testing.T) {
	// Testing tag not existing
	target := Object{Sha: "tag"}
	tag := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusNotFound).
		JSON(tag)

	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagUnauthorized(t *testing.T) {
	// Testing a 401
	body := release{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusUnauthorized).
		JSON(body)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagSuccessful(t *testing.T) {
	// Testing 201 response
	body := release{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusCreated).
		JSON(body)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagSuccessfulWithHostOverride(t *testing.T) {
	// Testing 201 response
	body := release{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.personal-github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.personal-github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusCreated).
		JSON(body)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "https://api.personal-github.com", "hello"}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagAlreadyExists(t *testing.T) {
	target := Object{Sha: "hash"}
	tag := Tag{Object: target}
	errorMessage := Error{Code: "already_exists"}
	response := BadResponse{Errors: []Error{errorMessage}}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK).
		JSON(tag)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusUnprocessableEntity).
		JSON(response)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "test", "hash", "", "hello"}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagOtherError(t *testing.T) {
	errorMessage := Error{Code: "blah"}
	response := BadResponse{Errors: []Error{errorMessage}}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusBadRequest).
		JSON(response)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherErrorResponse(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}
	assertTest.False(repo.CreateTag())
}
