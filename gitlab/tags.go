package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	urllib "net/url"
	"os"
)

// Gitlab interface for bitbucket
type Gitlab interface {
	create() bool
	validate() bool
}

// RepoProperties properties for repo
type RepoProperties struct {
	Token string
	Repo  string
	Tag   string
	Hash  string
	Host  string
	Body  string
}

// Commit Structure of bitbucket tag target
type Commit struct {
	ID string `json:"id"`
}

// Tag Structure of bitbucket tag response
type Tag struct {
	Commit Commit `json:"commit"`
}

// Release Object
type Release struct {
	Description string `json:"description"`
}

// BadResponse format for 400 http response body
type BadResponse struct {
	Message string `json:"message"`
}

//ValidateTag checks a tag does not exist or has the same hash
func (r *RepoProperties) ValidateTag() bool {
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := false
	url := ""
	if r.Host == "" {
		url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/tags/%s", urllib.QueryEscape(r.Repo), r.Tag)
	} else {
		url = fmt.Sprintf("%s/api/v4/projects/%s/repository/tags/%s", r.Host, urllib.QueryEscape(r.Repo), r.Tag)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error validate tag request")
	}
	request.Header.Set("PRIVATE-TOKEN", r.Token)
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
		if r.Hash == res.Commit.ID {
			validTag = true
		}
	}
	return validTag
}

// CreateTag creates a Gitlab tag
func (r *RepoProperties) CreateTag() bool {
	createTag := false
	if r.ValidateTag() {
		url := ""
		if r.Host == "" {
			url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/tags", urllib.QueryEscape(r.Repo))
		} else {
			url = fmt.Sprintf("%s/api/v4/projects/%s/repository/tags", r.Host, urllib.QueryEscape(r.Repo))
		}

		request, err := http.NewRequest("POST", url, nil)
		q := request.URL.Query()
		q.Add("tag_name", r.Tag)
		q.Add("ref", r.Hash)
		request.URL.RawQuery = q.Encode()
		if err != nil {
			fmt.Println("Error creating tag request", err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("PRIVATE-TOKEN", r.Token)
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
			// Create release notes with tag
			createTag = r.createRelease()
		}

		if resp.StatusCode == http.StatusBadRequest {
			res := BadResponse{}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body of error response")
			}
			err = json.Unmarshal(body, &res)
			if res.Message == fmt.Sprintf("Tag %s already exists", r.Tag) {
				createTag = r.createRelease()
			}
		}
	}
	return createTag
}

func (r *RepoProperties) createRelease() bool {
	createdRelease := false
	release := ""
	if r.Host == "" {
		release = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/tags/%s/release", urllib.QueryEscape(r.Repo), r.Tag)
	} else {
		release = fmt.Sprintf("%s/api/v4/projects/%s/repository/tags/%s/release", r.Host, urllib.QueryEscape(r.Repo), r.Tag)
	}
	body := Release{r.Body}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error marshalling object:", err)
	}
	request, err := http.NewRequest("POST", release, bytes.NewBuffer(jsonBody))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("PRIVATE-TOKEN", r.Token)
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

	if resp.StatusCode == http.StatusConflict {
		res := BadResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body of error response")
		}
		err = json.Unmarshal(body, &res)
		if res.Message == "Release already exists" {
			createdRelease = true
		}
	}

	if resp.StatusCode == http.StatusCreated {
		createdRelease = true
	}
	return createdRelease
}
