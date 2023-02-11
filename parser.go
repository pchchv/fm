package main

import "fmt"

type expr interface {
	String() string
	eval(app *app, args []string)
}

type callExpr struct {
	name  string
	args  []string
	count int
}

func (e *callExpr) String() string {
	return fmt.Sprintf("%s -- %s", e.name, e.args)
}
