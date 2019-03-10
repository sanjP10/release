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

func TestValidate_SetFlags(t *testing.T) {
	type args struct {
		f *flag.FlagSet
	}
	tests := []struct {
		name string
		v    *Validate
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.SetFlags(tt.args.f)
		})
	}
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
	type args struct {
		v *Validate
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkValidateFlags(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkValidateFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}
