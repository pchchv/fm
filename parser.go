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

type mapExpr struct {
	keys string
	expr expr
}

type cmapExpr struct {
	key  string
	expr expr
}

type cmdExpr struct {
	name string
	expr expr
}

func (e *callExpr) String() string {
	return fmt.Sprintf("%s -- %s", e.name, e.args)
}

func (e *setExpr) String() string {
	return fmt.Sprintf("set %s %s", e.opt, e.val)
}

func (e *mapExpr) String() string {
	return fmt.Sprintf("map %s %s", e.keys, e.expr)
}

func (e *cmapExpr) String() string {
	return fmt.Sprintf("cmap %s %s", e.key, e.expr)
}

func (e *cmdExpr) String() string {
	return fmt.Sprintf("cmd %s %s", e.name, e.expr)
}
