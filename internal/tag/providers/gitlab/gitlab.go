package gitlab

import (
	"bitbucket.org/cloudreach/release/internal/tag"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	urllib "net/url"
	"os"
)

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

// Properties properties for repo
type Properties struct {
	tag.RepoProperties
	Body string
}

//ValidateTag checks a tag does not exist or has the same hash
func (r *Properties) ValidateTag() tag.ValidTagState {
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := tag.ValidTagState{TagDoesntExist: false, TagExistsWithProvidedHash: false}
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
	if request == nil {
		_, err := os.Stderr.WriteString("Error creating request\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
		return validTag
	}
	request.Header.Set("PRIVATE-TOKEN", r.Password)
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
			validTag.TagExistsWithProvidedHash = true
		}
	}
	return validTag
}

// CreateTag creates a Gitlab tag
func (r *Properties) CreateTag() bool {
	createTag := false
	validTagState := r.ValidateTag()
	if validTagState.TagExistsWithProvidedHash {
		createTag = true
	} else if validTagState.TagDoesntExist {
		url := ""
		if r.Host == "" {
			url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/tags", urllib.QueryEscape(r.Repo))
		} else {
			url = fmt.Sprintf("%s/api/v4/projects/%s/repository/tags", r.Host, urllib.QueryEscape(r.Repo))
		}

		request, err := http.NewRequest("POST", url, nil)
		if request == nil {
			_, err := os.Stderr.WriteString("Error creating request\n")
			if err != nil {
				panic("Cannot write to stderr")
			}
			return false
		}
		q := request.URL.Query()
		q.Add("tag_name", r.Tag)
		q.Add("ref", r.Hash)
		request.URL.RawQuery = q.Encode()
		if err != nil {
			fmt.Println("Error creating tag request", err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("PRIVATE-TOKEN", r.Password)
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
			_, errorWriting := os.Stderr.WriteString(res.Message)
			if errorWriting != nil {
				panic("Cannot write to stderr")
			}
		}
	}
	return createTag
}

func (r *Properties) createRelease() bool {
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
	if request == nil {
		_, err := os.Stderr.WriteString("Error creating request\n")
		if err != nil {
			panic("Cannot write to stderr")
		}
		return false
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("PRIVATE-TOKEN", r.Password)
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
		_, errorWriting := os.Stderr.WriteString(res.Message)
		if errorWriting != nil {
			panic("Cannot write to stderr")
		}
	}

	if resp.StatusCode == http.StatusCreated {
		createdRelease = true
	}
	return createdRelease
}
