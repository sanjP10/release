package bitbucket

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestValidateTagNotExisting(t *testing.T) {
	// Testing tag not existing
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.True(repo.ValidateTag())
}

func TestValidateTagUnauthorized(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusUnauthorized)
	assertTest := assert.New(t)
	// Testing a 403
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.False(repo.ValidateTag())
}

func TestValidateTagExistingSameHash(t *testing.T) {
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusOK).
		JSON(tag)

	assertTest := assert.New(t)
	// Testing 200 response and hash is the same
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.True(repo.ValidateTag())
}

func TestValidateTagExistingMismatchHash(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response but hash is not the same
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusOK).
		JSON(tag)
	repo := RepoProperties{"username", "password", "repo", "tag", "not_hash", ""}
	assertTest.False(repo.ValidateTag())
}

func TestValidateTagOtherError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing a 403
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.False(repo.ValidateTag())
}

func TestCreateTagNotFound(t *testing.T) {
	// Testing tag not existing
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusNotFound).
		JSON(tag)

	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagUnauthorized(t *testing.T) {
	// Testing a 401
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusUnauthorized).
		JSON(tag)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagSuccessful(t *testing.T) {
	// Testing 201 response
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusCreated).
		JSON(tag)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagSuccessfulWithHostOverride(t *testing.T) {
	// Testing 201 response
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.personal-bitbucket.com").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.personal-bitbucket.com").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusCreated).
		JSON(tag)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", "api.personal-bitbucket.com"}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagAlreadyExists(t *testing.T) {
	errorMessage := Error{Message: "tag \"test\" already exists"}
	response := BadResponse{Type: "error", Error: errorMessage}
	defer gock.Off() // Flush pending mocks after test execution
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/test").
		Reply(http.StatusOK).
		JSON(tag)
	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusBadRequest).
		JSON(response)
	assertTest := assert.New(t)
	repo := RepoProperties{"username", "password", "repo", "test", "hash", ""}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagOtherError(t *testing.T) {
	errorMessage := Error{Message: "something went wrong"}
	response := BadResponse{Type: "error", Error: errorMessage}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusOK)
	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusBadRequest).
		JSON(response)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherErrorResponse(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusOK)
	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := RepoProperties{"username", "password", "repo", "tag", "hash", ""}
	assertTest.False(repo.CreateTag())
}
