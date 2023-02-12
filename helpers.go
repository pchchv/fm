package main

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

var (
	reAltKey  = regexp.MustCompile(`<a-(.)>`)
	reWord    = regexp.MustCompile(`(\pL|\pN)+`)
	reWordBeg = regexp.MustCompile(`([^\pL\pN]|^)(\pL|\pN)`)
	reWordEnd = regexp.MustCompile(`(\pL|\pN)([^\pL\pN]|$)`)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runeSliceWidth(rs []rune) (w int) {
	for _, r := range rs {
		w += runewidth.RuneWidth(r)
	}
	return
}

func runeSliceWidthRange(rs []rune, beg, end int) []rune {
	curr := 0
	b := 0

	for i, r := range rs {
		w := runewidth.RuneWidth(r)

		switch {
		case curr == beg:
			b = i
		case curr < beg && curr+w > beg:
			b = i + 1
		case curr == end:
			return rs[b:i]
		case curr > end:
			return rs[b : i-1]
		}
		curr += w
	}
	return nil
}

// humanize converts the size in bytes to human-readable form using metric suffixes (e.g., 1K = 1000).
// For values less than 10, the first significant digit is shown, otherwise it is hidden.
// Numbers are always rounded down.
// This should suit most people.
func humanize(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / 1000
	for _, s := range suffix {
		if curr < 10 {
			return fmt.Sprintf("%.1f%s", curr-0.0499, s)
		} else if curr < 1000 {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= 1000
	}

	return ""
}

// naturalLess compares two strings for a natural sort that takes into account
// the values of the numbers in the strings.
// For example, "2" is less than "10", and similarly "foo2bar" is less than "foo10bar",
// but "bar2bar" is greater than "foo10bar".
func naturalLess(s1, s2 string) bool {
	lo1, lo2, hi1, hi2 := 0, 0, 0, 0
	for {
		if hi1 >= len(s1) {
			return hi2 != len(s2)
		}

		if hi2 >= len(s2) {
			return false
		}

		isDigit1 := isDigit(s1[hi1])
		isDigit2 := isDigit(s2[hi2])

		for lo1 = hi1; hi1 < len(s1) && isDigit(s1[hi1]) == isDigit1; hi1++ {
		}

		for lo2 = hi2; hi2 < len(s2) && isDigit(s2[hi2]) == isDigit2; hi2++ {
		}

		if s1[lo1:hi1] == s2[lo2:hi2] {
			continue
		}

		if isDigit1 && isDigit2 {
			num1, err1 := strconv.Atoi(s1[lo1:hi1])
			num2, err2 := strconv.Atoi(s2[lo2:hi2])

			if err1 == nil && err2 == nil {
				return num1 < num2
			}
		}

		return s1[lo1:hi1] < s2[lo2:hi2]
	}
}

func isRoot(name string) bool {
	return filepath.Dir(name) == name
}

// escape is used to escape spaces and special characters with backspaces in a given string.
func escape(s string) string {
	buf := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsSpace(r) || r == '\\' || r == ';' || r == '#' {
			buf = append(buf, '\\')
		}
		buf = append(buf, r)
	}
	return string(buf)
}

// unescape is used to remove backspaces,
// which are used to escape spaces and special characters in a given string.
func unescape(s string) string {
	esc := false
	buf := make([]rune, 0, len(s))
	for _, r := range s {
		if esc {
			if !unicode.IsSpace(r) && r != '\\' && r != ';' && r != '#' {
				buf = append(buf, '\\')
			}
			buf = append(buf, r)
			esc = false
			continue
		}
		if r == '\\' {
			esc = true
			continue
		}
		esc = false
		buf = append(buf, r)
	}

	if esc {
		buf = append(buf, '\\')
	}

	return string(buf)
}
