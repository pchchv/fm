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

type setExpr struct {
	opt string
	val string
}

func (e *callExpr) String() string {
	return fmt.Sprintf("%s -- %s", e.name, e.args)
}

func (e *setExpr) String() string {
	return fmt.Sprintf("set %s %s", e.opt, e.val)
}
