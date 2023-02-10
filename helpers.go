package main

import (
	"bufio"
	"fmt"
	"io"
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

// ReadPairs reads pairs of lines, separated by spaces, on each line.
// Single or double quotes can be used to escape spaces.
// Hash characters can be used to add a comment before the end of a line.
// Leading and trailing spaces are truncated. Blank lines are skipped.
func readPairs(r io.Reader) ([][]string, error) {
	var pairs [][]string

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		squote, dquote := false, false
		for i := 0; i < len(line); i++ {
			if line[i] == '\'' && !dquote {
				squote = !squote
			} else if line[i] == '"' && !squote {
				dquote = !dquote
			}
			if !squote && !dquote && line[i] == '#' {
				line = line[:i]
				break
			}
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		squote, dquote = false, false
		pair := strings.FieldsFunc(line, func(r rune) bool {
			if r == '\'' && !dquote {
				squote = !squote
			} else if r == '"' && !squote {
				dquote = !dquote
			}
			return !squote && !dquote && unicode.IsSpace(r)
		})

		if len(pair) != 2 {
			return nil, fmt.Errorf("expected pair but found: %s", s.Text())
		}

		for i := 0; i < len(pair); i++ {
			squote, dquote = false, false
			buf := make([]rune, 0, len(pair[i]))
			for _, r := range pair[i] {
				if r == '\'' && !dquote {
					squote = !squote
					continue
				}
				if r == '"' && !squote {
					dquote = !dquote
					continue
				}
				buf = append(buf, r)
			}
			pair[i] = string(buf)
		}

		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func replaceTilde(s string) string {
	if strings.HasPrefix(s, "~") {
		s = strings.Replace(s, "~", genUser.HomeDir, 1)
	}
	return s
}
