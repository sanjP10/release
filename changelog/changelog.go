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

// GetVersions retrieves versions from changelog where lines are formatted as with prefix of ## and numbers
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

// ValidateVersionSemantics takes the desired and previous versions and ensures that the desired is larger than previous
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

// RetrieveChanges gets all the changes from the log file between the desired and previous version lines
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

// ConvertToDesiredTag changes the markdown version line into a version tag, by removing markdown notation and spaces
func (c *Properties) ConvertToDesiredTag() string {
	markdownRegex := regexp.MustCompile("##\\s*")
	return markdownRegex.ReplaceAllString(c.desired, "")
}

// ReadChangelogAsString reads a file and returns it as a string
func ReadChangelogAsString(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	changelog := string(dat)
	return changelog, err
}

func convertVersionToFloat(version string) float64 {
	// convert string ## x.x.x to a float
	r := regexp.MustCompile("##\\s*|\\.")
	versionAsFloatString := r.ReplaceAllString(version, "")
	versionFloat, _ := strconv.ParseFloat(versionAsFloatString, 64)
	return versionFloat
}
