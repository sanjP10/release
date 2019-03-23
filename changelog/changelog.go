package changelog

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Bitbucket interface for bitbucket
type Changelog interface {
	create() bool
	validate() bool
}

type ChangelogProperties struct {
	previous string
	desired  string
	changes  string
}

func (c *ChangelogProperties) getVersions(changelog string) {
	changelogNum := regexp.MustCompile("##\\s\\d.+")
	matches := changelogNum.FindAllString(changelog, -1)
	if len(matches) > 1 {
		c.desired = matches[0]
		c.previous = matches[1]
	} else if len(matches) == 1 {
		c.desired = matches[0]
	}
}

func validateVersionSemantics(desiredVersion string, previousVersion string) bool {
	previousVersionFloat, _ := convertVersionToFloat(previousVersion)
	desiredVersionFloat, _ := convertVersionToFloat(desiredVersion)
	return desiredVersionFloat > previousVersionFloat
}

func convertVersionToFloat(version string) (float64, error) {
	r := regexp.MustCompile("##\\s+|\\.")
	versionAsFloatString := r.ReplaceAllString(version, "")
	versionFloat, err := strconv.ParseFloat(versionAsFloatString, 64)
	return versionFloat, err
}

func retrieveChanges(previousVersion string, desiredVersion string, changelog string) {
	scanner := bufio.NewScanner(strings.NewReader(changelog))
	startRecording := false
	var changes []string
	for scanner.Scan() {
		if scanner.Text() == desiredVersion {
			startRecording = true
			continue
		}
		if previousVersion != "" && scanner.Text() == previousVersion {
			startRecording = false
			break
		}
		if startRecording && scanner.Text() != "" {
			changes = append(changes, scanner.Text())
		}
	}
	fmt.Println(strings.Join(changes, "\n"))
}

func readFileAsString(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	changelog := string(dat)
	return changelog, err
}

func getChanges(filename string) {
	changelog, _ := readFileAsString(filename)
	matches := getVersions(changelog)
	if len(matches) > 1 {
		validScemantics := validateVersionSemantics(matches[0], matches[1])
		if validScemantics {
			retrieveChanges(matches[1], matches[0], changelog)
		}
	} else if len(matches) == 1 {
		retrieveChanges("", matches[0], changelog)
	} else {
		// error
	}
}
