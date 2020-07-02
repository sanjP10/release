package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate_Name(t *testing.T) {
	validate := &Validate{}
	assertTest := assert.New(t)
	assertTest.Equal(validate.Name(), "validate")
}

func TestValidate_Synopsis(t *testing.T) {
	validate := &Validate{}
	assertTest := assert.New(t)
	assertTest.Equal(validate.Synopsis(), "Validates tag and release for repo to be created.")
}

func TestValidate_Usage(t *testing.T) {
	validate := &Validate{}
	assertTest := assert.New(t)
	var expected = "Validates tag and release for repo to be created.\n"
	assertTest.Equal(validate.Usage(), expected)
}

func Test_checkValidateFlagsGit(t *testing.T) {
	validateCmd := &Validate{}
	errors := checkValidateFlags(validateCmd)
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

	validateCmd.username = "testuser"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-password required",
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	validateCmd.password = "password"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	validateCmd.repo = "repo"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-changelog required",
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	validateCmd.changelog = "changelog"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-hash required",
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	validateCmd.hash = "hash"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-email required",
		"-origin required"}
	assertTest.Equal(errors, expected)

	validateCmd.email = "an-email@abc.com"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-origin required"}
	assertTest.Equal(errors, expected)

	validateCmd.origin = "http://an-origin.com/repo.git"
	validCreate := checkValidateFlags(validateCmd)
	assertTest.Empty(validCreate)
}

func Test_checkValidateFlagsGitSSH(t *testing.T) {
	validateCmd := &Validate{}
	errors := checkValidateFlags(validateCmd)
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

	validateCmd.ssh = "ssh-file"
	validateCmd.origin = "http://an-origin.com/repo.git"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-changelog required",
		"-hash required",
		"-email required"}
	assertTest.Equal(errors, expected)

	validateCmd.changelog = "changelog"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-hash required",
		"-email required"}
	assertTest.Equal(errors, expected)

	validateCmd.hash = "hash"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-email required"}
	assertTest.Equal(errors, expected)

	validateCmd.email = "an-email@abc.com"
	validCreate := checkValidateFlags(validateCmd)
	assertTest.Empty(validCreate)
}

func TestValidate_SetFlags(t *testing.T) {
	validateCmd := &Validate{}
	assertTest := assert.New(t)
	validateCmd.provider = "svn"
	errors := checkValidateFlags(validateCmd)
	expected := []string{
		"-username required",
		"-password required",
		"-repo required",
		"-changelog required",
		"-hash required",
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

	for _, provider := range providers {
		validateCmd.provider = provider
		errors := checkValidateFlags(validateCmd)
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

	validateCmd.username = "testuser"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-password required",
		"-repo required",
		"-changelog required",
		"-hash required",
	}
	assertTest.Equal(errors, expected)

	validateCmd.password = "password"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-repo required",
		"-changelog required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	validateCmd.repo = "repo"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-changelog required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	validateCmd.changelog = "changelog"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-hash required"}
	assertTest.Equal(errors, expected)

	validateCmd.hash = "hash"
	assertTest.Empty(checkValidateFlags(validateCmd))
}

func Test_ValidateCheckFlag_Gitlab(t *testing.T) {
	validate := &Validate{}
	validate.password = "token"
	validate.provider = "gitlab"
	validate.repo = "repo"
	validate.hash = "hash"
	validate.changelog = "file"
	assertTest := assert.New(t)
	errors := checkValidateFlags(validate)
	assertTest.Empty(errors)
}
