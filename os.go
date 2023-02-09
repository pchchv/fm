package main

import (
	"os/exec"
	"os/user"

	"golang.org/x/sys/unix"
)

var (
	genDefaultShell      = "sh"
	genDefaultShellFlag  = "-c"
	genDefaultSocketProt = "unix"
	genDefaultSocketPath string

	genUser        *user.User
	genConfigPaths []string
	genColorsPaths []string
	genIconsPaths  []string
	genFilesPath   string
	genMarksPath   string
	genTagsPath    string
	genHistoryPath string
)

func detachedCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &unix.SysProcAttr{Setsid: true}
	return cmd
}
