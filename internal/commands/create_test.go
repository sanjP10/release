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
		"-username required",
		"-password required",
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required",
	}
	assertTest := assert.New(t)
	assertTest.Equal(errors, expected)

	createCmd.username = "testuser"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-password required",
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	createCmd.password = "password"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	createCmd.repo = "repo"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	createCmd.changelog = "changelog"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	createCmd.hash = "hash"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	createCmd.email = "an-email@abc.com"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-origin required"}
	assertTest.Equal(errors, expected)

	createCmd.origin = "http://an-origin.com/repo.git"
	validCreate := checkCreateFlags(createCmd)
	assertTest.Empty(validCreate)
}

func Test_checkCreateFlagsGitSSH(t *testing.T) {
	createCmd := &Create{}
	errors := checkCreateFlags(createCmd)
	expected := []string{
		"-username required",
		"-password required",
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required",
	}
	assertTest := assert.New(t)
	assertTest.Equal(errors, expected)

	createCmd.ssh = "ssh-file"
	createCmd.origin = "http://an-origin.com/repo.git"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-changelog required",
		"-hash required",
		"-email required"}
	assertTest.Equal(errors, expected)

	createCmd.changelog = "changelog"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-hash required",
		"-email required"}
	assertTest.Equal(errors, expected)

	createCmd.hash = "hash"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-email required"}
	assertTest.Equal(errors, expected)

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
		"-username required",
		"-password required",
		"-repo required",
		"-changelog required",
		"-hash required",
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

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
		assertTest.Equal(errors, expected)
	}

	createCmd.username = "testuser"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-password required",
		"-repo required",
		"-changelog required",
		"-hash required",
	}
	assertTest.Equal(errors, expected)

	createCmd.password = "password"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-repo required",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	createCmd.repo = "repo"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-changelog required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	createCmd.changelog = "changelog"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-hash required"}
	assertTest.Equal(errors, expected)

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
