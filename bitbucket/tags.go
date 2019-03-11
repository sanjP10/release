package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type commitId struct {
	hash string
}

type target struct {
	Hash commitId `json:"hash"`
}

type Tag struct {
	name   string
	Target target `json:"target"`
}

func ValidateTag(username string, password string, repo string, tag string, hash string) bool {
	// Check tag exists, if 404 gd, 403 auth error, 200 exists and check hash is the same
	validTag := false
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags/%s", repo, tag)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error validate tag request")
	}
	request.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error validate tag request")
	}

	if resp.StatusCode == 404 {
		validTag = true
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
		if hash == res.Target.Hash.hash {
			validTag = true
		}
	}
	return validTag
}

func CreateTag(username string, password string, repo string, tag string, hash string) bool {
	createdTag := false
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/refs/tags", repo)
	commit := commitId{hash: hash}
	targetObj := target{Hash: commit}
	body := Tag{name: tag, Target: targetObj}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error marshalling object:", err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating tag request")
	}
	request.Header.Add("Content-Type", "application/json")
	request.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error creating tag")
	}
	if resp.StatusCode == 201 {
		createdTag = true
	} else {
		fmt.Println("Error creating tag")
	}
	return createdTag
}
