package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InvalidProvider(t *testing.T) {
	assertTest := assert.New(t)
	assertTest.False(ValidProvider("svn"))
}

func Test_ValidProvider(t *testing.T) {
	provider := "GITHUB"
	assertTest := assert.New(t)
	assertTest.True(ValidProvider(provider))

	provider = "Github"
	assertTest.True(ValidProvider(provider))

	provider = "gitHub"
	assertTest.True(ValidProvider(provider))

	provider = "GITLAB"
	assertTest.True(ValidProvider(provider))

	provider = "GitLab"
	assertTest.True(ValidProvider(provider))

	provider = "gitlab"
	assertTest.True(ValidProvider(provider))

	provider = "BITBUCKET"
	assertTest.True(ValidProvider(provider))

	provider = "BitBucket"
	assertTest.True(ValidProvider(provider))

	provider = "bitbucket"
	assertTest.True(ValidProvider(provider))
}
