package cmd

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
	assertTest.Equal(create.Synopsis(), "create release for bitbucket repo.")
}

func TestCreate_Usage(t *testing.T) {
	create := &Create{}
	assertTest := assert.New(t)
	expected := `create [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>] [-host <host> (optional)]:
  creates tag against bitbucket repo
`
	assertTest.Equal(create.Usage(), expected)
}

func Test_checkCreateFlags(t *testing.T) {
	createCmd := &Create{}
	errors := checkCreateFlags(createCmd)
	expected := []string{
		"-username required",
		"-password required",
		"-repo required",
		"-tag required",
		"-hash required"}
	assertTest := assert.New(t)
	assertTest.Equal(errors, expected)

	createCmd.username = "testuser"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-password required",
		"-repo required",
		"-tag required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	createCmd.password = "password"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-repo required",
		"-tag required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	createCmd.repo = "repo"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-tag required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	createCmd.tag = "tag"
	errors = checkCreateFlags(createCmd)
	expected = []string{
		"-hash required"}
	assertTest.Equal(errors, expected)

	createCmd.hash = "hash"
	validCreate := checkCreateFlags(createCmd)
	assertTest.Empty(validCreate)
}
