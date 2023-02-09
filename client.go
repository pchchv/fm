package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

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
