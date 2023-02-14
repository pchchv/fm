package main

import (
	"os"
	"reflect"
	"testing"
)

func TestIsRoot(t *testing.T) {
	sep := string(os.PathSeparator)
	if !isRoot(sep) {
		t.Errorf(`"%s" is root`, sep)
	}

	paths := []string{
		"",
		"~",
		"foo",
		"foo/bar",
		"foo/bar",
		"/home",
		"/home/user",
	}

	for _, p := range paths {
		if isRoot(p) {
			t.Errorf("'%s' is not root", p)
		}
	}
}

func TestRuneSliceWidth(t *testing.T) {
	tests := []struct {
		rs  []rune
		exp int
	}{
		{[]rune{'a', 'b'}, 2},
		{[]rune{'ı', 'ş'}, 2},
		{[]rune{'世', '界'}, 4},
		{[]rune{'世', 'a', '界', 'ı'}, 6},
	}

	for _, test := range tests {
		if got := runeSliceWidth(test.rs); got != test.exp {
			t.Errorf("at input '%v' expected '%d' but got '%d'", test.rs, test.exp, got)
		}
	}
}

func TestRuneSliceWidthRange(t *testing.T) {
	tests := []struct {
		rs  []rune
		beg int
		end int
		exp []rune
	}{
		{[]rune{'a', 'b', 'c', 'd'}, 1, 3, []rune{'b', 'c'}},
		{[]rune{'a', 'ı', 'b', 'ş'}, 1, 3, []rune{'ı', 'b'}},
		{[]rune{'世', '界', '世', '界'}, 2, 6, []rune{'界', '世'}},
		{[]rune{'世', '界', '世', '界'}, 3, 6, []rune{'世'}},
		{[]rune{'世', '界', '世', '界'}, 2, 5, []rune{'界'}},
		{[]rune{'世', '界', '世', '界'}, 3, 5, []rune{}},
		{[]rune{'世', 'a', '界', 'ı'}, 2, 5, []rune{'a', '界'}},
		{[]rune{'世', 'a', '界', 'ı'}, 2, 4, []rune{'a'}},
		{[]rune{'世', 'a', '界', 'ı'}, 3, 5, []rune{'界'}},
		{[]rune{'世', 'a', '界', 'ı'}, 3, 4, []rune{}},
	}

	for _, test := range tests {
		if got := runeSliceWidthRange(test.rs, test.beg, test.end); !reflect.DeepEqual(got, test.exp) {
			t.Errorf("at input '%v' expected '%v' but got '%v'", test.rs, test.rs, got)
		}
	}
}
