package bitbucket

import (
	"encoding/json"
	"github.com/karupanerura/go-mock-http-response"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func mockClient(statusCode int, headers map[string]string, body []byte) *http.Client {
	return mockhttp.NewResponseMock(statusCode, headers, body).MakeClient()
}

func TestValidateTagNotExisting(t *testing.T) {
	// Testing tag not existing
	assertTest := assert.New(t)
	client := mockClient(http.StatusNotFound, nil, nil)
	assertTest.True(ValidateTag("username",
		"password",
		"repo", "tag", "hash", *client))
}

func TestValidateTagUnauthorized(t *testing.T) {
	assertTest := assert.New(t)
	// Testing a 403
	client := mockClient(http.StatusUnauthorized, nil, nil)
	assertTest.False(ValidateTag("username",
		"password",
		"repo", "tag", "hash", *client))
}

func TestValidateTagExistingSameHash(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response and hash is the same
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	jsonTag, _ := json.Marshal(tag)
	client := mockClient(http.StatusOK, nil, jsonTag)
	assertTest.True(ValidateTag("username",
		"password",
		"repo", "tag", "hash", *client))
}

func TestValidateTagExistingMismatchHash(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response but hash is not the same
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	jsonTag, _ := json.Marshal(tag)
	client := mockClient(http.StatusOK, nil, jsonTag)
	assertTest.False(ValidateTag("username",
		"password",
		"repo", "tag", "not_hash", *client))
}

func TestCreateTag(t *testing.T) {
	// Testing tag not existing
	assertTest := assert.New(t)
	client := mockClient(http.StatusNotFound, nil, nil)
	assertTest.False(CreateTag("username",
		"password",
		"repo", "tag", "hash", *client))
}

func TestCreateTagUnauthorized(t *testing.T) {
	assertTest := assert.New(t)
	// Testing a 403
	client := mockClient(http.StatusUnauthorized, nil, nil)
	assertTest.False(CreateTag("username",
		"password",
		"repo", "tag", "hash", *client))
}

func TestCreateTagSuccessful(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 201 response
	target := Target{Hash: "hash"}
	tag := Tag{Name: "tag", Target: target}
	jsonTag, _ := json.Marshal(tag)
	client := mockClient(http.StatusCreated, nil, jsonTag)
	assertTest.True(CreateTag("username",
		"password",
		"repo", "tag", "hash", *client))
}

func TestCreateTagAlreadyExists(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	errorMessage := Error{Message: "tag \"test\" already exists"}
	tag := BadResponse{Type: "error", Error: errorMessage}
	jsonTag, _ := json.Marshal(tag)
	client := mockClient(http.StatusBadRequest, nil, jsonTag)
	assertTest.True(CreateTag("username",
		"password",
		"repo", "test", "hash", *client))
}

func TestCreateTagOtherError(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	errorMessage := Error{Message: "something went wrong"}
	tag := BadResponse{Type: "error", Error: errorMessage}
	jsonTag, _ := json.Marshal(tag)
	client := mockClient(http.StatusBadRequest, nil, jsonTag)
	assertTest.False(CreateTag("username",
		"password",
		"repo", "test", "hash", *client))
}
