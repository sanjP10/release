package cmd

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
	assertTest.Equal(validate.Synopsis(), "validates release version to be created.")
}

func TestValidate_Usage(t *testing.T) {
	validate := &Validate{}
	assertTest := assert.New(t)
	var expected = `validate [-username <username>] [-password <password/token>] [-repo <repo>] [-changelog <changelog md file>] [-provider <github/gitlab/bitbucket>] [-host <host> (optional)]:
  validates tag against bitbucket repo
`
	assertTest.Equal(validate.Usage(), expected)
}

func Test_checkValidateFlags(t *testing.T) {
	validateCmd := &Validate{}
	errors := checkValidateFlags(validateCmd)
	expected := []string{
		"-username required",
		"-password required",
		"-repo required",
		"-changelog required",
		"-hash required",
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest := assert.New(t)
	assertTest.Equal(errors, expected)

	validateCmd.username = "testuser"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-password required",
		"-repo required",
		"-changelog required",
		"-hash required",
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

	validateCmd.password = "password"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-repo required",
		"-changelog required",
		"-hash required",
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

	validateCmd.repo = "repo"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-changelog required",
		"-hash required",
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

	validateCmd.changelog = "changelog"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-hash required",
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

	validateCmd.hash = "hash"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

	validateCmd.provider = "svn"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-provider required, valid values are github, gitlab, bitbucket"}
	assertTest.Equal(errors, expected)

	for _, provider := range providers {
		validateCmd.provider = provider
		validCreate := checkValidateFlags(validateCmd)
		assertTest.Empty(validCreate)
	}
}
