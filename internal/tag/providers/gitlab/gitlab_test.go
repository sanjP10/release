package gitlab

import (
	"bitbucket.org/cloudreach/release/internal/tag/interfaces"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestValidateTagNotExisting_Gitlab(t *testing.T) {
	// Testing tag not existing
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusNotFound)
	assertTest := assert.New(t)
	repo := GitlabProperties{RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.True(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagUnauthorized_Gitlab(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusUnauthorized)
	assertTest := assert.New(t)
	// Testing a 403
	repo := GitlabProperties{RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingSameHash_Gitlab(t *testing.T) {
	target := Commit{ID: "hash"}
	tag := GitlabTag{Commit: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusOK).
		JSON(tag)

	assertTest := assert.New(t)
	// Testing 200 response and hash is the same
	repo := GitlabProperties{RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.True(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingMismatchHash_Gitlab(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response but hash is not the same
	target := Commit{ID: "hash"}
	tag := GitlabTag{Commit: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusOK).
		JSON(tag)
	repo := GitlabProperties{RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "not_hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagOtherError_Gitlab(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing a 403
	repo := GitlabProperties{RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestCreateTagNotFound_Gitlab(t *testing.T) {
	// Testing tag not existing
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "tag").
		MatchParam("ref", "hash").
		Reply(http.StatusNotFound)

	assertTest := assert.New(t)
	repo := GitlabProperties{RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagUnauthorized_Gitlab(t *testing.T) {
	// Testing a 401
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "tag").
		MatchParam("ref", "hash").
		Reply(http.StatusUnauthorized)
	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagSuccessful_Gitlab(t *testing.T) {
	// Testing 201 response
	body := GitlabRelease{Description: "hello"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "tag").
		MatchParam("ref", "hash").
		Reply(http.StatusCreated)

	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags/tag/release").
		Reply(http.StatusCreated).
		JSON(body)

	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagSuccessfulWithHostOverride_Gitlab(t *testing.T) {
	// Testing 201 response
	body := GitlabRelease{"hello"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("http://personal-gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusNotFound)

	gock.New("http://personal-gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "tag").
		MatchParam("ref", "hash").
		Reply(http.StatusCreated)

	gock.New("http://personal-gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags/tag/release").
		Reply(http.StatusCreated).
		JSON(body)
	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", "http://personal-gitlab.com"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagAndReleaseAlreadyExists_Gitlab(t *testing.T) {
	target := Commit{ID: "hash"}
	tag := GitlabTag{Commit: target}

	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/test").
		Reply(http.StatusOK).
		JSON(tag)
	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "test", "hash", ""}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagAndErrors_Gitlab(t *testing.T) {
	target := Commit{ID: "hash"}
	tag := GitlabTag{Commit: target}
	response := GitlabBadResponse{"Tag test already exists"}
	body := GitlabRelease{"hello"}
	releaseResponse := GitlabBadResponse{"GitlabRelease already exists"}

	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/test").
		Reply(http.StatusNotFound).
		JSON(tag)
	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "test").
		MatchParam("ref", "hash").
		Reply(http.StatusCreated).
		JSON(response)
	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags/test/release").
		JSON(body).
		Reply(http.StatusConflict).
		JSON(releaseResponse)
	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "test", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagAndReleaseFails(t *testing.T) {
	target := Commit{ID: "hash"}
	tag := GitlabTag{Commit: target}
	response := GitlabBadResponse{"Tag test already exists"}
	body := GitlabRelease{"hello"}

	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/test").
		Reply(http.StatusNotFound).
		JSON(tag)
	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "test").
		MatchParam("ref", "hash").
		Reply(http.StatusCreated).
		JSON(response)
	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags/test/release").
		JSON(body).
		Reply(http.StatusNotFound)
	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "test", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherError_Gitlab(t *testing.T) {
	response := GitlabBadResponse{"blah"}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusNotFound)
	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "tag").
		MatchParam("ref", "hash").
		Reply(http.StatusBadRequest).
		JSON(response)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherErrorResponse_Gitlab(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://gitlab.com/").
		Get("api/v4/projects/org/repo/repository/tags/tag").
		Reply(http.StatusNotFound)
	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags").
		MatchParam("tag_name", "tag").
		MatchParam("ref", "hash").
		Reply(http.StatusInternalServerError)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateReleaseNotFound(t *testing.T) {
	// Testing release not existing
	body := GitlabRelease{Description: "hello"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags/tag/release").
		Reply(http.StatusNotFound).
		JSON(body)

	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.False(repo.createRelease())
}

func TestCreateReleaseUnauthorized(t *testing.T) {
	// Testing a 401
	// Testing release not existing
	body := GitlabRelease{Description: "hello"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags/tag/release").
		Reply(http.StatusUnauthorized).
		JSON(body)

	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.False(repo.createRelease())
}

func TestCreateRelease(t *testing.T) {
	// Testing release not existing
	body := GitlabRelease{"hello"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://gitlab.com/").
		Post("api/v4/projects/org/repo/repository/tags/tag/release").
		JSON(body).
		Reply(http.StatusCreated).
		JSON(body)

	assertTest := assert.New(t)
	repo := GitlabProperties{Body: "hello", RepoProperties: interfaces.RepoProperties{"token", "org/repo", "tag", "hash", ""}}
	assertTest.True(repo.createRelease())
}
