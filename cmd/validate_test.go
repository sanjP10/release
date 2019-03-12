package cmd

import (
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/google/subcommands"
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
	var expected = `validate [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>]:
  validates tag against bitbucket repo
`
	assertTest.Equal(validate.Usage(), expected)
}

func TestValidate_Execute(t *testing.T) {
	type args struct {
		in0 context.Context
		f   *flag.FlagSet
		in2 []interface{}
	}
	tests := []struct {
		name string
		v    *Validate
		args args
		want subcommands.ExitStatus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Execute(tt.args.in0, tt.args.f, tt.args.in2...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkValidateFlags(t *testing.T) {
	validateCmd := &Validate{}
	errors := checkValidateFlags(validateCmd)
	expected := []string{
		"-username required",
		"-password required",
		"-repo required",
		"-tag required",
		"-hash required"}
	assertTest := assert.New(t)
	assertTest.Equal(errors, expected)

	validateCmd.username = "testuser"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-password required",
		"-repo required",
		"-tag required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	validateCmd.password = "password"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-repo required",
		"-tag required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	validateCmd.repo = "repo"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-tag required",
		"-hash required"}
	assertTest.Equal(errors, expected)

	validateCmd.tag = "tag"
	errors = checkValidateFlags(validateCmd)
	expected = []string{
		"-hash required"}
	assertTest.Equal(errors, expected)

	validateCmd.hash = "hash"
	validCreate := checkValidateFlags(validateCmd)
	assertTest.Empty(validCreate)
}
