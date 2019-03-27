package changelog

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestValidateVersionSemanticsNoPrevious(t *testing.T) {
	// No previous
	changelog := &Properties{previous: "", desired: "## 1.0.0"}
	assertTest := assert.New(t)
	assertTest.True(changelog.ValidateVersionSemantics())
}

func TestValidateVersionSemantics(t *testing.T) {
	// checking when previous is below the desired
	changelog := &Properties{previous: "## 0.0.1.0", desired: "## 1.0.0.0"}
	assertTest := assert.New(t)
	assertTest.True(changelog.ValidateVersionSemantics())
}

func TestValidateVersionSemanticsInvalid(t *testing.T) {
	// checking when previous is above the desired
	changelog := &Properties{previous: "## 2.0.1", desired: "## 1.10.0"}
	assertTest := assert.New(t)
	assertTest.False(changelog.ValidateVersionSemantics())
}

func TestGetVersion(t *testing.T) {
	// Checking version string formatted as '## major.minor.patch' has markdown removed
	assertTest := assert.New(t)
	actual := getVersion("## 1.0.0")
	assertTest.Equal("1.0.0", actual)
	actual = getVersion("## 1.0.0     ")
	assertTest.Equal("1.0.0", actual)
	actual = getVersion("##    1.0.0     ")
	assertTest.Equal("1.0.0", actual)
}

func TestConvertVersionToFloatNoSpace(t *testing.T) {
	// Checking version string formatted as '## major.minor.patch' has markdown removed
	assertTest := assert.New(t)
	actual := getVersion("##1.0.0")
	assertTest.Equal("1.0.0", actual)
}

func TestReadChangelogAsStringNoFile(t *testing.T) {
	// file doesn't exist
	assertTest := assert.New(t)
	file, err := ReadChangelogAsString("blah.md")
	assertTest.Empty(file)
	assertTest.Error(err)
}

func TestReadChangelogAsString(t *testing.T) {
	// checking file is returned as string
	assertTest := assert.New(t)
	dat, _ := ioutil.ReadFile("../fixtures/Changelog.md")
	expected := string(dat)
	file, err := ReadChangelogAsString("../fixtures/Changelog.md")
	assertTest.Equal(file, expected)
	assertTest.Empty(err)
}

func TestGetVersionsFirst(t *testing.T) {
	assertTest := assert.New(t)
	file, _ := ReadChangelogAsString("../fixtures/FirstChangelog.md")
	changelog := &Properties{}
	changelog.GetVersions(file)
	assertTest.Equal("", changelog.previous)
	assertTest.Equal("##0.0.0", changelog.desired)
}

func TestGetVersions(t *testing.T) {
	assertTest := assert.New(t)
	file, _ := ReadChangelogAsString("../fixtures/Changelog.md")
	changelog := &Properties{}
	changelog.GetVersions(file)
	assertTest.Equal("##1.0.0", changelog.previous)
	assertTest.Equal("##    1.1.0", changelog.desired)
}

func TestRetrieveChanges(t *testing.T) {
	assertTest := assert.New(t)
	file, _ := ReadChangelogAsString("../fixtures/Changelog.md")
	changelog := &Properties{}
	changelog.GetVersions(file)
	changelog.RetrieveChanges(file)
	assertTest.Equal(`### Updated
* An update happened`, changelog.Changes)
}

func TestFirstRetrieveChanges(t *testing.T) {
	assertTest := assert.New(t)
	file, _ := ReadChangelogAsString("../fixtures/FirstChangelog.md")
	changelog := &Properties{}
	changelog.GetVersions(file)
	changelog.RetrieveChanges(file)
	assertTest.Equal(`### Added
* Initial release`, changelog.Changes)
}

func TestConvertToDesiredTag(t *testing.T) {
	assertTest := assert.New(t)
	changelog := &Properties{desired: "##    1.1.0"}
	assertTest.Equal("1.1.0", changelog.ConvertToDesiredTag())

	changelog.desired = "##1.1.0"
	assertTest.Equal("1.1.0", changelog.ConvertToDesiredTag())
}
