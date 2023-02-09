package main

import "golang.org/x/sys/windows"

var (
	genDefaultShell      = "cmd"
	genDefaultShellFlag  = "/c"
	genDefaultSocketProt = "tcp"
	genDefaultSocketPath = "127.0.0.1:12345"

	genUser        *user.User
	genConfigPaths []string
	genColorsPaths []string
	genIconsPaths  []string
	genFilesPath   string
	genTagsPath    string
	genMarksPath   string
	genHistoryPath string
)

func detachedCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &windows.SysProcAttr{CreationFlags: 8}
	return cmd
}

func setUserUmask() {}
