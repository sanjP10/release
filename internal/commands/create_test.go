package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate_Name(t *testing.T) {
	create := &Create{}
	assertTest := assert.New(t)
	assertTest.Equal(create.Name(), "create")
}

func TestCreate_Synopsis(t *testing.T) {
	create := &Create{}
	assertTest := assert.New(t)
	assertTest.Equal(create.Synopsis(), "Creates tag and release for repo.")
}

func TestCreate_Usage(t *testing.T) {
	create := &Create{}
	assertTest := assert.New(t)
	expected := "Creates tag and release for repo.\n"
	assertTest.Equal(create.Usage(), expected)
}

func Test_checkCreateFlagsGit(t *testing.T) {
	createCmd := &Create{}
	errors := checkCreateFlags(createCmd)
	expected := []string{
		"-origin required",
		"-username or -ssh required, for CodeCommit or GCP Source repositories both are required",
		"-password required",
		"-email required",
		"-changelog required",
		"-hash required"}
	assertTest := assert.New(t)
	assertTest.Equal(expected, errors)

	createCmd.username = "tester"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-origin required",
		"-password required",
		"-email required",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.password = "password"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-origin required",
		"-email required",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.repo = "repo"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-origin required",
		"-email required",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.changelog = "changelog"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-origin required",
		"-email required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.hash = "hash"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-origin required",
		"-email required"}
	assertTest.Equal(expected, errors)

	createCmd.email = "an-email@abc.com"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-origin required"}
	assertTest.Equal(expected, errors)

	createCmd.origin = "https://an-origin.com/repo.git"
	validCreate := checkCreateFlags(createCmd)
	assertTest.Empty(validCreate)
}

func Test_checkCreateFlagsGitSSH(t *testing.T) {
	createCmd := &Create{}
	errors := checkCreateFlags(createCmd)
	expected := []string{
		"-origin required",
		"-username or -ssh required, for CodeCommit or GCP Source repositories both are required",
		"-password required",
		"-email required",
		"-changelog required",
		"-hash required"}
	assertTest := assert.New(t)
	assertTest.Equal(expected, errors)

	createCmd.ssh = "ssh-file"
	createCmd.origin = "https://an-origin.com/repo.git"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-email required",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.changelog = "changelog"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-email required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.hash = "hash"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-email required"}
	assertTest.Equal(expected, errors)

	createCmd.email = "an-email@abc.com"
	validCreate := checkCreateFlags(createCmd)
	assertTest.Empty(validCreate)
}

func TestCreate_SetFlags(t *testing.T) {
	createCmd := &Create{}
	assertTest := assert.New(t)
	createCmd.provider = "svn"
	errors := checkCreateFlags(createCmd)
	expected := []string{
		"-provider valid values are github, gitlab, bitbucket",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	for _, provider := range providers {
		createCmd.provider = provider
		errors := checkCreateFlags(createCmd)
		expected := []string{
			"-username required",
			"-password required",
			"-repo required",
			"-changelog required",
			"-hash required"}
		if provider == "gitlab" {
			expected = append(expected[:0], expected[1:]...)
		}
		assertTest.Equal(expected, errors)
	}

	createCmd.username = "tester"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-password required",
		"-repo required",
		"-changelog required",
		"-hash required",
	}
	assertTest.Equal(expected, errors)

	createCmd.password = "password"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-repo required",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.repo = "repo"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-changelog required",
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.changelog = "changelog"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-hash required"}
	assertTest.Equal(expected, errors)

	createCmd.hash = "hash"
	assertTest.Empty(checkCreateFlags(createCmd))
}

func Test_CreateCheckFlag_Gitlab(t *testing.T) {
	create := &Create{}
	create.password = "token"
	create.provider = "gitlab"
	create.repo = "repo"
	create.hash = "hash"
	create.changelog = "file"
	assertTest := assert.New(t)
	errors := checkCreateFlags(create)
	assertTest.Empty(errors)
}
