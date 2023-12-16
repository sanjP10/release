package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sanjP10/release/internal/tag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

// BadResponse structure of 400 response
type BadResponse struct {
	Type  string `json:"type"`
	Error Error  `json:"error"`
}

// Error structure of error message response
type Error struct {
	Message string `json:"message"`
}

// Properties for repo
type Properties struct {
	tag.RepoProperties
	Username string
	Repo     string
	Host     string
}

// ServerTag response
type ServerTag struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	Hash            string `json:"hash"`
}

// ServerTagBody tag body for request
type ServerTagBody struct {
	Name       string `json:"name"`
	StartPoint string `json:"startPoint"`
	Message    string `json:"message"`
}

// ValidateTag checks a tag does not exist or has the same hash
func (r *Properties) ValidateTag() tag.ValidTagState {
	isCloud := true // Using bitbucket cloud offering otherwise self-hosted
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := tag.ValidTagState{TagDoesntExist: false, TagExistsWithProvidedHash: false}
	url := ""
	if r.Host == "" {
		url = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags/%s", r.Repo, r.Tag)
	} else {
		isCloud = false
		repoDetails := strings.Split(r.Repo, "/")
		url = fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/tags/%s", r.Host, repoDetails[0], repoDetails[1], r.Tag)
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
	switch resp.StatusCode {
	case http.StatusNotFound:
		validTag.TagDoesntExist = true
	case http.StatusUnauthorized:
		_, err := os.Stderr.WriteString("Unauthorised, please check credentials\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
	case http.StatusOK:
		validTag.TagExistsWithProvidedHash = checkResponse(resp, r.Hash, isCloud)
	}
	return validTag
}

// CreateTag creates a bitbucket tag
func (r *Properties) CreateTag() bool {
	createTag := false
	validTagState := r.ValidateTag()
	if validTagState.TagExistsWithProvidedHash {
		createTag = true
	} else if validTagState.TagDoesntExist {
		isCloud := true // Using bitbucket cloud offering otherwise self-hosted
		url := ""
		if r.Host == "" {
			url = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags", r.Repo)
		} else {
			isCloud = false
			repoDetails := strings.Split(r.Repo, "/")
			url = fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/tags", r.Host, repoDetails[0], repoDetails[1])
		}

		jsonBody := createBody(r, isCloud)

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

		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusNotFound:
			_, err := os.Stderr.WriteString("Unauthorised, please check credentials\n")
			if err != nil {
				panic("Cannot write to stderr")
			}
		case http.StatusOK, http.StatusCreated:
			createTag = true
		case http.StatusBadRequest:
			res := BadResponse{}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body of error response")
			}
			err = json.Unmarshal(body, &res)
			if err != nil {
				fmt.Println("Error unmarshalling response")
			}
			_, errorWriting := os.Stderr.WriteString(res.Error.Message)
			if errorWriting != nil {
				panic("Cannot write to stderr")
			}
		}

	}
	return createTag
}

func createBody(r *Properties, isCloud bool) []byte {
	var jsonBody []byte
	var err error
	if isCloud {
		target := Target{r.Hash}
		body := &Tag{Name: r.Tag, Target: target}

		jsonBody, err = json.Marshal(body)
		if err != nil {
			fmt.Println("error marshalling object:", err)
		}
	} else {
		body := &ServerTagBody{Name: r.Tag, StartPoint: r.Hash, Message: r.Body}
		jsonBody, err = json.Marshal(body)
		if err != nil {
			fmt.Println("error marshalling object:", err)
		}
	}
	return jsonBody
}

func checkResponse(resp *http.Response, hash string, isCloud bool) bool {
	existsWithProvidedHash := false
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body of tag response")
	}
	if isCloud {
		res := Tag{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println("Error unmarshalling body")
		}
		if hash == res.Target.Hash {
			existsWithProvidedHash = true
		}
	} else {
		res := ServerTag{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println("Error unmarshalling body")
		}
		if hash == res.LatestCommit {
			existsWithProvidedHash = true
		}
	}
	return existsWithProvidedHash
}
