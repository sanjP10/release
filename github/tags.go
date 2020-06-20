package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Github interface for github
type Github interface {
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
	Body     string
}

// Object Structure of gitlab tag target
type Object struct {
	Sha string `json:"sha"`
}

// Tag Structure of github tag response
type Tag struct {
	Object Object `json:"object"`
}

// release struct format required for github release api
type release struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	Prerelease      bool   `json:"prerelease"`
}

// Error structure of error message response
type Error struct {
	Code string `json:"code"`
}

// BadResponse format for 400 http response body
type BadResponse struct {
	Message string  `json:"message"`
	Errors  []Error `json:"errors"`
}

//ValidateTag checks a tag does not exist or has the same hash
func (r *RepoProperties) ValidateTag() bool {
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := false
	url := ""
	if r.Host == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/git/refs/tags/%s", r.Repo, r.Tag)
	} else {
		url = fmt.Sprintf("%s/repos/%s/git/refs/tags/%s", r.Host, r.Repo, r.Tag)
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
		if r.Hash == res.Object.Sha {
			validTag = true
		}
	}
	return validTag
}

// CreateTag creates a github tag
func (r *RepoProperties) CreateTag() bool {
	createTag := false
	if r.ValidateTag() {
		url := ""
		if r.Host == "" {
			url = fmt.Sprintf("https://api.github.com/repos/%s/releases", r.Repo)
		} else {
			url = fmt.Sprintf("%s/repos/%s/releases", r.Host, r.Repo)
		}

		body := release{Name: r.Tag, TagName: r.Tag, Body: r.Body, Draft: false, Prerelease: false, TargetCommitish: r.Hash}

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

		if resp.StatusCode == http.StatusUnprocessableEntity {
			res := BadResponse{}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body of error response")
			}
			err = json.Unmarshal(body, &res)
			if res.Errors[0].Code == fmt.Sprintf("already_exists") {
				createTag = true
			}
		}
	}
	return createTag
}
