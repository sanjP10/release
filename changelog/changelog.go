package changelog

import (
	"bufio"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Changelog interface for changelog.md file
type Changelog interface {
	getChanges()
	RetrieveChanges()
	GetVersions()
	ConvertToDesiredTag()
}

// Properties for changelog
type Properties struct {
	previous string
	desired  string
	Changes  string
}

func (c *Properties) GetVersions(changelog string) {
	changelogNumRegex := regexp.MustCompile("##\\s*\\d.+")
	matches := changelogNumRegex.FindAllString(changelog, -1)
	if len(matches) > 1 {
		c.desired = matches[0]
		c.previous = matches[1]
	} else if len(matches) == 1 {
		c.desired = matches[0]
	}
}

func (c *Properties) ValidateVersionSemantics() bool {
	valid := false
	if c.previous == "" {
		valid = true
	} else {
		previousVersionFloat := convertVersionToFloat(c.previous)
		desiredVersionFloat := convertVersionToFloat(c.desired)
		valid = desiredVersionFloat > previousVersionFloat
	}
	return valid
}

func (c *Properties) RetrieveChanges(changelog string) {
	scanner := bufio.NewScanner(strings.NewReader(changelog))
	startRecording := false
	var changes []string
	for scanner.Scan() {
		if scanner.Text() == c.desired {
			startRecording = true
			continue
		}
		if c.previous != "" && scanner.Text() == c.previous {
			startRecording = false
			break
		}
		if startRecording && scanner.Text() != "" {
			changes = append(changes, scanner.Text())
		}
	}
	c.Changes = strings.Join(changes, "\n")
}


func (c *Properties) ConvertToDesiredTag() string {
	markdownRegex := regexp.MustCompile("##\\s*")
	return markdownRegex.ReplaceAllString(c.desired, "")
}

func ReadChangelogAsString(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	changelog := string(dat)
	return changelog, err
}

func convertVersionToFloat(version string) float64 {
	r := regexp.MustCompile("##\\s*|\\.")
	versionAsFloatString := r.ReplaceAllString(version, "")
	versionFloat, _ := strconv.ParseFloat(versionAsFloatString, 64)
	return versionFloat
}
