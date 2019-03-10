package cmd

import (
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/google/subcommands"
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
	expected := `create [-username <username>] [-password <password/token>] [-repo <repo>] [-tag <tag>]:
  creates tag against bitbucket repo
`
	assertTest.Equal(create.Usage(), expected)
}

func TestCreate_SetFlags(t *testing.T) {
	type args struct {
		f *flag.FlagSet
	}
	tests := []struct {
		name string
		c    *Create
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetFlags(tt.args.f)
		})
	}
}

func TestCreate_Execute(t *testing.T) {
	type args struct {
		in0 context.Context
		f   *flag.FlagSet
		in2 []interface{}
	}
	tests := []struct {
		name string
		c    *Create
		args args
		want subcommands.ExitStatus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Execute(tt.args.in0, tt.args.f, tt.args.in2...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkCreateFlags(t *testing.T) {
	type args struct {
		c *Create
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
			if got := checkCreateFlags(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkCreateFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}
