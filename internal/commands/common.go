package commands

import (
	"strings"
)

var providers = [...]string{"github", "gitlab", "bitbucket"}

// ValidProvider Check provider from cli is supported
func ValidProvider(str string) bool {
	for _, a := range providers {
		if a == strings.ToLower(str) {
			return true
		}
	}
	return false
}
