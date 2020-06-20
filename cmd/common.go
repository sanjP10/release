package cmd

import (
	"strings"
)

var providers = [...]string{"github", "gitlab", "bitbucket"}

func ValidProvider(str string) bool {
	for _, a := range providers {
		if a == strings.ToLower(str) {
			return true
		}
	}
	return false
}
