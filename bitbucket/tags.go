package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Bitbucket interface for bitbucket
type Bitbucket interface {
	create() bool
	validate() bool
}

// RepoProperties properties for repo
type RepoProperties struct {
	Username string
	Password string
	Repo     string
	Tag      string
	Hash     string
	Host     string
}

// Target Structure of bitbucket tag target
type Target struct {
	Hash string `json:"hash"`
}

// Tag Structure of bitbucket tag response
type Tag struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}

// BadResponse structure of 400 response
type BadResponse struct {
	Type  string `json:"type"`
	Error Error  `json:"error"`
}

// Error structure of error message response
type Error struct {
	Message string `json:"message"`
}

//ValidateTag checks a tag does not exist or has the same hash
func (r *RepoProperties) ValidateTag() bool {
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := false
	url := ""
	if r.Host == "" {
		url = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags/%s", r.Repo, r.Tag)
	} else {
		url = fmt.Sprintf("https://%s/2.0/repositories/%s/refs/tags/%s", r.Host, r.Repo, r.Tag)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error validate tag request")
	}
	request.SetBasicAuth(r.Username, r.Password)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error validate tag request")
	}

	if resp.StatusCode == http.StatusNotFound {
		validTag = true
	}
	if resp.StatusCode == http.StatusUnauthorized {
		_, err := os.Stderr.WriteString("Unauthorised, please check credentials\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
	}
	if resp.StatusCode == http.StatusOK {
		res := Tag{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body of tag response")
		}
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println("Error unmarshalling body")
		}
		if r.Hash == res.Target.Hash {
			validTag = true
		}
	}
	return validTag
}

// CreateTag creates a bitbucket tag
func (r *RepoProperties) CreateTag() bool {
	createTag := false
	if r.ValidateTag() {
		url := ""
		if r.Host == "" {
			url = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags", r.Repo)
		} else {
			url = fmt.Sprintf("https://%s/2.0/repositories/%s/refs/tags", r.Host, r.Repo)
		}

		target := Target{r.Hash}
		body := &Tag{Name: r.Tag, Target: target}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			fmt.Println("error marshalling object:", err)
		}
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("Error creating tag request", err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.SetBasicAuth(r.Username, r.Password)
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			fmt.Println("Error creating tag", err)
		}
		if resp.StatusCode == http.StatusUnauthorized {
			_, err := os.Stderr.WriteString("Unauthorised, please check credentials\n")
			if err != nil {
				panic("Cannot write to stderr")
			}
		}
		if resp.StatusCode == http.StatusNotFound {
			_, err := os.Stderr.WriteString("Repo not found\n")
			if err != nil {
				panic("Cannot write to stderr")
			}
		}

		if resp.StatusCode == http.StatusCreated {
			createTag = true
		}

		if resp.StatusCode == http.StatusBadRequest {
			res := BadResponse{}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body of error response")
			}
			err = json.Unmarshal(body, &res)
			if res.Error.Message == fmt.Sprintf("tag \"%s\" already exists", r.Tag) {
				createTag = true
			} else {
				_, errorWriting := os.Stderr.WriteString(res.Error.Message)
				if errorWriting != nil {
					panic("Cannot write to stderr")
				}
			}
		}
	}
	return createTag
}
