package bitbucket

import (
	"bitbucket.org/cloudreach/release/interfaces"
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
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.True(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagUnauthorized(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusUnauthorized)
	assertTest := assert.New(t)
	// Testing a 403
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingSameHash(t *testing.T) {
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusOK).
		JSON(tag)

	assertTest := assert.New(t)
	// Testing 200 response and hash is the same
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.True(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingMismatchHash(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response but hash is not the same
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusOK).
		JSON(tag)
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "not_hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagOtherError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing a 403
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestCreateTagNotFound(t *testing.T) {
	// Testing tag not existing
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusNotFound).
		JSON(tag)

	assertTest := assert.New(t)
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagUnauthorized(t *testing.T) {
	// Testing a 401
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusUnauthorized).
		JSON(tag)
	assertTest := assert.New(t)
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagSuccessful(t *testing.T) {
	// Testing 201 response
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusCreated).
		JSON(tag)
	assertTest := assert.New(t)
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagSuccessfulWithHostOverride(t *testing.T) {
	// Testing 201 response
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.personal-bitbucket.com").
		Get("/2.0/repositories/repo/refs/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://api.personal-bitbucket.com").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusCreated).
		JSON(tag)
	assertTest := assert.New(t)
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", "https://api.personal-bitbucket.com"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagAlreadyExists(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/test").
		Reply(http.StatusOK).
		JSON(tag)
	assertTest := assert.New(t)
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "test", "hash", ""}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagError(t *testing.T) {
	errorMessage := BitbucketError{Message: "tag \"test\" already exists"}
	response := BitbucketBadResponse{Type: "error", Error: errorMessage}
	defer gock.Off() // Flush pending mocks after test execution
	target := Target{Hash: "hash"}
	tag := BitbucketTag{Name: "tag", Target: target}
	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/repo/refs/tags/test").
		Reply(http.StatusNotFound).
		JSON(tag)
	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/repo/refs/tags").
		Reply(http.StatusBadRequest).
		JSON(response)
	assertTest := assert.New(t)
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "test", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherError(t *testing.T) {
	errorMessage := BitbucketError{Message: "something went wrong"}
	response := BitbucketBadResponse{Type: "error", Error: errorMessage}
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
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
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
	repo := BitbucketProperties{Username: "username", RepoProperties: interfaces.RepoProperties{"password", "repo", "tag", "hash", ""}}
	assertTest.False(repo.CreateTag())
}
