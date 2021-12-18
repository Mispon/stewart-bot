package utils

import "strings"

// HasAnyOf looking for any of words in string
func HasAnyOf(str string, words []string) bool {
	for _, w := range words {
		if strings.Contains(str, w) {
			return true
		}
	}
	return false
}
