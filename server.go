package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/pchchv/golog"
)

var (
	genListener net.Listener
	genConnList = make(map[int]net.Conn)
	genQuitChan = make(chan struct{}, 1)
)

func serve() {
	if genLogPath != "" {
		f, err := os.OpenFile(genLogPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	golog.Info("hi!")

	if genSocketProt == "unix" {
		setUserUmask()
	}

	l, err := net.Listen(genSocketProt, genSocketPath)
	if err != nil {
		golog.Info("listening socket: %s", err)
		return
	}
	defer l.Close()

	genListener = l

	listen(l)
}

func listen(l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			select {
			case <-genQuitChan:
				golog.Info("bye!")
				return
			default:
				golog.Info("accepting connection: %s", err)
			}
		}
		go handleConn(c)
	}
}

func echoerr(c net.Conn, msg string) {
	fmt.Fprintln(c, msg)
	golog.Info(msg)
}

func echoerrf(c net.Conn, format string, a ...interface{}) {
	echoerr(c, fmt.Sprintf(format, a...))
}

func handleConn(c net.Conn) {
	s := bufio.NewScanner(c)

Loop:
	for s.Scan() {
		golog.Info("listen: %s", s.Text())
		word, rest := splitWord(s.Text())
		switch word {
		case "conn":
			if rest != "" {
				word2, _ := splitWord(rest)
				id, err := strconv.Atoi(word2)
				if err != nil {
					echoerr(c, "listen: conn: client id should be a number")
				} else {
					genConnList[id] = c
				}
			} else {
				echoerr(c, "listen: conn: requires a client id")
			}
		case "drop":
			if rest != "" {
				word2, _ := splitWord(rest)
				id, err := strconv.Atoi(word2)
				if err != nil {
					echoerr(c, "listen: drop: client id should be a number")
				} else {
					delete(genConnList, id)
				}
			} else {
				echoerr(c, "listen: drop: requires a client id")
			}
		case "send":
			if rest != "" {
				word2, rest2 := splitWord(rest)
				id, err := strconv.Atoi(word2)
				if err != nil {
					for _, c := range genConnList {
						fmt.Fprintln(c, rest)
					}
				} else {
					if c2, ok := genConnList[id]; ok {
						fmt.Fprintln(c2, rest2)
					} else {
						echoerr(c, "listen: send: no such client id is connected")
					}
				}
			}
		case "quit":
			if len(genConnList) == 0 {
				genQuitChan <- struct{}{}
				genListener.Close()
				break Loop
			}
		case "quit!":
			genQuitChan <- struct{}{}
			for _, c := range genConnList {
				fmt.Fprintln(c, "echo server is quitting...")
				c.Close()
			}
			genListener.Close()
			break Loop
		default:
			echoerrf(c, "listen: unexpected command: %s", word)
		}
	}

	if s.Err() != nil {
		echoerrf(c, "listening: %s", s.Err())
	}

	c.Close()
}
