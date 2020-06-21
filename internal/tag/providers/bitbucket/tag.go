package bitbucket

import (
	"bitbucket.org/cloudreach/release/internal/tag/interfaces"
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

// BitbucketTag Structure of bitbucket tag response
type BitbucketTag struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}

// BitbucketBadResponse structure of 400 response
type BitbucketBadResponse struct {
	Type  string         `json:"type"`
	Error BitbucketError `json:"error"`
}

// BitbucketError structure of error message response
type BitbucketError struct {
	Message string `json:"message"`
}

// BitbucketProperties properties for repo
type BitbucketProperties struct {
	interfaces.RepoProperties
	Username string
}

//ValidateTag checks a tag does not exist or has the same hash
func (r *BitbucketProperties) ValidateTag() interfaces.ValidTagState {
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := interfaces.ValidTagState{TagDoesntExist: false, TagExistsWithProvidedHash: false}
	url := ""
	if r.Host == "" {
		url = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags/%s", r.Repo, r.Tag)
	} else {
		url = fmt.Sprintf("%s/2.0/repositories/%s/refs/tags/%s", r.Host, r.Repo, r.Tag)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error validate tag request")
	}
	if request == nil {
		_, err := os.Stderr.WriteString("Error creating request\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
		return validTag
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
		return validTag
	}
	if resp.StatusCode == http.StatusNotFound {
		validTag.TagDoesntExist = true
	}
	if resp.StatusCode == http.StatusUnauthorized {
		_, err := os.Stderr.WriteString("Unauthorised, please check credentials\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
	}
	if resp.StatusCode == http.StatusOK {
		res := BitbucketTag{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body of tag response")
		}
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println("Error unmarshalling body")
		}
		if r.Hash == res.Target.Hash {
			validTag.TagExistsWithProvidedHash = true
		}
	}
	return validTag
}

// CreateTag creates a bitbucket tag
func (r *BitbucketProperties) CreateTag() bool {
	createTag := false
	validTagState := r.ValidateTag()
	if validTagState.TagExistsWithProvidedHash {
		createTag = true
	} else if validTagState.TagDoesntExist {
		url := ""
		if r.Host == "" {
			url = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags", r.Repo)
		} else {
			url = fmt.Sprintf("%s/2.0/repositories/%s/refs/tags", r.Host, r.Repo)
		}

		target := Target{r.Hash}
		body := &BitbucketTag{Name: r.Tag, Target: target}

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
			res := BitbucketBadResponse{}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body of error response")
			}
			err = json.Unmarshal(body, &res)
			_, errorWriting := os.Stderr.WriteString(res.Error.Message)
			if errorWriting != nil {
				panic("Cannot write to stderr")
			}
		}
	}
	return createTag
}
