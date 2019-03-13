package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Target Structure of bitbucket tag target
type Target struct {
	Hash string `json:"hash"`
}

// Tag Structure of bitbucket tag response
type Tag struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}

//ValidateTag checks a tag does not exist or has the same hash
func ValidateTag(username string, password string, repo string, tag string, hash string, client http.Client) bool {
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := false
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags/%s", repo, tag)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error validate tag request")
	}
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error validate tag request")
	}

	if resp.StatusCode == 404 {
		validTag = true
	}
	if resp.StatusCode == 401 {
		_, err := os.Stderr.WriteString("Unauthorised, please check credentials\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
	}
	if resp.StatusCode == 200 {
		res := Tag{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body of tag response")
		}
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println("Error unmarshalling body")
		}
		if hash == res.Target.Hash {
			validTag = true
		}
	}
	return validTag
}

// CreateTag creates a bitbucket tag
func CreateTag(username string, password string, repo string, tag string, hash string, client http.Client) bool {
	createdTag := false
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags", repo)
	target := Target{hash}
	body := &Tag{Name: tag, Target: target}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error marshalling object:", err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating tag request", err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error creating tag", err)
	}
	if resp.StatusCode == 401 {
		_, err := os.Stderr.WriteString("Unauthorised, please check credentials\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
	}
	if resp.StatusCode == 404 {
		_, err := os.Stderr.WriteString("Repo not found\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
	}

	if resp.StatusCode == 201 {
		createdTag = true
	}
	return createdTag
}
