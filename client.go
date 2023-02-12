package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pchchv/golog"
)

func run() {
	var screen tcell.Screen
	var err error

	if screen, err = tcell.NewScreen(); err != nil {
		golog.Fatal("creating screen: %s", err)
	} else if err = screen.Init(); err != nil {
		golog.Fatal("initializing screen: %s", err)
	}
	if genOpts.mouse {
		screen.EnableMouse()
	}

	if genLogPath != "" {
		f, err := os.OpenFile(genLogPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	golog.Info("hi!")

	ui := newUI(screen)
	nav := newNav(ui.wins[0].h)
	app := newApp(ui, nav)

	if err := nav.sync(); err != nil {
		app.ui.echoerrf("sync: %s", err)
	}

	if err := app.nav.readMarks(); err != nil {
		app.ui.echoerrf("reading marks file: %s", err)
	}

	if err := app.nav.readTags(); err != nil {
		app.ui.echoerrf("reading tags file: %s", err)
	}

	if err := app.readHistory(); err != nil {
		app.ui.echoerrf("reading history file: %s", err)
	}

	app.loop()

	app.ui.screen.Fini()
}

func remote(cmd string) error {
	c, err := net.Dial(genSocketProt, genSocketPath)
	if err != nil {
		return fmt.Errorf("dialing to send server: %s", err)
	}

	fmt.Fprintln(c, cmd)

	// the standard net.Conn interface does not include the CloseWrite method,
	// but net.UnixConn and net.TCPConn implement it,
	// so the following should be safe as long as no other connection types are used.
	// CloseWrite is needed to notify the server that this is not a persistent connection
	// and should be closed after a response.
	if v, ok := c.(interface {
		CloseWrite() error
	}); ok {
		v.CloseWrite()
	}

	io.Copy(os.Stdout, c)

	c.Close()

	return nil
}

func readExpr() <-chan expr {
	ch := make(chan expr)

	go func() {
		duration := 1 * time.Second

		c, err := net.Dial(genSocketProt, genSocketPath)
		for err != nil {
			golog.Info("connecting server: %s", err)
			time.Sleep(duration)
			duration *= 2
			c, err = net.Dial(genSocketProt, genSocketPath)
		}

		fmt.Fprintf(c, "conn %d\n", genClientID)

		ch <- &callExpr{"sync", nil, 1}

		s := bufio.NewScanner(c)
		for s.Scan() {
			golog.Info("recv: %s", s.Text())
			p := newParser(strings.NewReader(s.Text()))
			if p.parse() {
				ch <- p.expr
			}
		}

		c.Close()
	}()

	return ch
}
