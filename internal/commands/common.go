package commands

import (
	"strings"
)

var providers = [...]string{"github", "gitlab", "bitbucket"}

// ValidProvider Check provider from cli is supported
func ValidProvider(str string) bool {
	isValid := false
	for _, a := range providers {
		if a == strings.ToLower(str) {
			isValid = true
			break
		}
	}
	return isValid
}
