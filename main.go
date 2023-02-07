package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	singleMode    bool
	clientID      int
	hostname      string
	lastDirPath   string
	selectionPath string
	socketProt    string
	socketPath    string
	logPath       string
	sSelect       string
	configPath    string
	flags         arrayFlag
	version       string
)

type arrayFlag []string

func (a *arrayFlag) Set(v string) error {
	*a = append(*a, v)
	return nil
}

func (a *arrayFlag) String() string {
	return strings.Join(*a, ", ")
}

func main() {
	flag.Usage = func() {
		f := flag.CommandLine.Output()
		fmt.Fprintln(f, "fm - Terminal file manager")
		fmt.Fprintln(f, "")
		fmt.Fprintf(f, "Usage:  %s [options] [cd-or-select-path]\n\n", os.Args[0])
		fmt.Fprintln(f, "  cd-or-select-path")
		fmt.Fprintln(f, "        set the initial dir or file selection to the given argument")
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "Options:")
		flag.PrintDefaults()
	}

	showDoc := flag.Bool("doc", false, "show documentation")
	showVersion := flag.Bool("version", false, "show version")
	serverMode := flag.Bool("server", false, "start server (automatic)")
	singleMode := flag.Bool("single", false, "start a client without server")
	remoteCmd := flag.String("remote", "", "send remote command to server")
	cpuprofile := flag.String("cpuprofile", "", "path to the file to write the CPU profile")
	memprofile := flag.String("memprofile", "", "path to the file to write the memory profile")
	flag.StringVar(&lastDirPath, "last-dir-path", "", "path to the file to write the last dir on exit (to use for cd)")
	flag.StringVar(&selectionPath, "selection-path", "", "path to the file to write selected files on open (to use as open file dialog)")
	flag.StringVar(&configPath, "config", "", "path to the config file (instead of the usual paths)")
	flag.Var(&flags, "command", "command to execute on client initialization")
	flag.StringVar(&logPath, "log", "", "path to the log file to write messages")

	flag.Parse()
}
