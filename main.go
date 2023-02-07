package main

import (
	"strings"
)

type arrayFlag []string

func (a *arrayFlag) Set(v string) error {
	*a = append(*a, v)
	return nil
}

func (a *arrayFlag) String() string {
	return strings.Join(*a, ", ")
}

func main() {}
