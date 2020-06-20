package tagging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Object Structure of gitlab tag target
type Object struct {
	Sha string `json:"sha"`
}

// Tag Structure of github tag response
type Tag struct {
	Object Object `json:"object"`
}

// GithubRelease struct format required for github release api
type GithubRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	Prerelease      bool   `json:"prerelease"`
}

// GithubError structure of error message response
type GithubError struct {
	Code string `json:"code"`
}

// GithubBadResponse format for 400 http response body
type GithubBadResponse struct {
	Message string        `json:"message"`
	Errors  []GithubError `json:"errors"`
}

//ValidateTag checks a tag does not exist or has the same hash
func (r *GithubProperties) ValidateTag() bool {
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
		fmt.Println("GithubError validate tag request")
	}
	if request == nil {
		_, err := os.Stderr.WriteString("Error creating request\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
		return false
	}
	request.SetBasicAuth(r.Username, r.Password)
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error validate tag request")
	}
	if resp == nil {
		_, err := os.Stderr.WriteString("Error getting response\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
		return false
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
			fmt.Println("GithubError unmarshalling body")
		}
		if r.Hash == res.Object.Sha {
			validTag = true
		}
	}
	return validTag
}

// CreateTag creates a github tag
func (r *GithubProperties) CreateTag() bool {
	createTag := false
	if r.ValidateTag() {
		url := ""
		if r.Host == "" {
			url = fmt.Sprintf("https://api.github.com/repos/%s/releases", r.Repo)
		} else {
			url = fmt.Sprintf("%s/repos/%s/releases", r.Host, r.Repo)
		}

		body := GithubRelease{Name: r.Tag, TagName: r.Tag, Body: r.Body, Draft: false, Prerelease: false, TargetCommitish: r.Hash}

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
			res := GithubBadResponse{}
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
