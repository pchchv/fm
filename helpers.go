package main

import (
	"strings"
	"unicode"
)

// splitWord splits the first word of a space-separated string from the rest.
// This is used to tokenize the string one word at a time without affecting the rest.
// Spaces to the left of the word and the rest of the word are trimmed.
func splitWord(s string) (word, rest string) {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	ind := len(s)

	for i, c := range s {
		if unicode.IsSpace(c) {
			ind = i
			break
		}
	}

	word = s[0:ind]
	rest = strings.TrimLeftFunc(s[ind:], unicode.IsSpace)
	return
}
