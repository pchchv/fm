package main

type expr interface {
	String() string
	eval(app *app, args []string)
}
