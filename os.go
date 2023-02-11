package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	envOpener = os.Getenv("OPENER")
	envEditor = os.Getenv("EDITOR")
	envPager  = os.Getenv("PAGER")
	envShell  = os.Getenv("SHELL")

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

func init() {
	if envOpener == "" {
		if runtime.GOOS == "darwin" {
			envOpener = "open"
		} else {
			envOpener = "xdg-open"
		}
	}

	if envEditor == "" {
		envEditor = "vi"
	}

	if envPager == "" {
		envPager = "less"
	}

	if envShell == "" {
		envShell = "sh"
	}

	u, err := user.Current()
	if err != nil {
		// When the user is not in /etc/passwd (for e.g. LDAP) and CGO_ENABLED=1 in go env,
		// the cgo implementation of user.Current() fails even when HOME and USER are set.

		log.Printf("user: %s", err)
		if os.Getenv("HOME") == "" {
			panic("$HOME variable is empty or not set")
		}
		if os.Getenv("USER") == "" {
			panic("$USER variable is empty or not set")
		}
		u = &user.User{
			Username: os.Getenv("USER"),
			HomeDir:  os.Getenv("HOME"),
		}
	}
	genUser = u

	config := os.Getenv("XDG_CONFIG_HOME")
	if config == "" {
		config = filepath.Join(genUser.HomeDir, ".config")
	}

	genConfigPaths = []string{
		filepath.Join("/etc", "lf", "lfrc"),
		filepath.Join(config, "lf", "lfrc"),
	}

	genColorsPaths = []string{
		filepath.Join("/etc", "lf", "colors"),
		filepath.Join(config, "lf", "colors"),
	}

	genIconsPaths = []string{
		filepath.Join("/etc", "lf", "icons"),
		filepath.Join(config, "lf", "icons"),
	}

	data := os.Getenv("XDG_DATA_HOME")
	if data == "" {
		data = filepath.Join(genUser.HomeDir, ".local", "share")
	}

	genFilesPath = filepath.Join(data, "lf", "files")
	genMarksPath = filepath.Join(data, "lf", "marks")
	genTagsPath = filepath.Join(data, "lf", "tags")
	genHistoryPath = filepath.Join(data, "lf", "history")

	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if runtime == "" {
		runtime = os.TempDir()
	}

	genDefaultSocketPath = filepath.Join(runtime, fmt.Sprintf("lf.%s.sock", genUser.Username))
}

func detachedCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &unix.SysProcAttr{Setsid: true}
	return cmd
}

func shellCommand(s string, args []string) *exec.Cmd {
	if len(genOpts.ifs) != 0 {
		s = fmt.Sprintf("IFS='%s'; %s", genOpts.ifs, s)
	}

	args = append([]string{genOpts.shellflag, s, "--"}, args...)

	args = append(genOpts.shellopts, args...)

	return exec.Command(genOpts.shell, args...)
}

func shellSetPG(cmd *exec.Cmd) {
	cmd.SysProcAttr = &unix.SysProcAttr{Setpgid: true}
}

func shellKill(cmd *exec.Cmd) error {
	pgid, err := unix.Getpgid(cmd.Process.Pid)
	if err == nil && cmd.Process.Pid == pgid {
		// kill the process group
		err = unix.Kill(-pgid, 15)
		if err == nil {
			return nil
		}
	}
	return cmd.Process.Kill()
}

func setDefaults() {
	genOpts.cmds["open"] = &execExpr{"&", `$OPENER "$f"`}
	genOpts.keys["e"] = &execExpr{"$", `$EDITOR "$f"`}
	genOpts.keys["i"] = &execExpr{"$", `$PAGER "$f"`}
	genOpts.keys["w"] = &execExpr{"$", "$SHELL"}

	genOpts.cmds["doc"] = &execExpr{"$", "lf -doc | $PAGER"}
	genOpts.keys["<f-1>"] = &callExpr{"doc", nil, 1}
}

func setUserUmask() {
	unix.Umask(0o077)
}

func isExecutable(f os.FileInfo) bool {
	return f.Mode()&0o111 != 0
}

func isHidden(f os.FileInfo, path string, hiddenfiles []string) bool {
	hidden := false
	for _, pattern := range hiddenfiles {
		matched := matchPattern(strings.TrimPrefix(pattern, "!"), f.Name(), path)
		if strings.HasPrefix(pattern, "!") && matched {
			hidden = false
		} else if matched {
			hidden = true
		}
	}
	return hidden
}

func userName(f os.FileInfo) string {
	if stat, ok := f.Sys().(*syscall.Stat_t); ok {
		if u, err := user.LookupId(fmt.Sprint(stat.Uid)); err == nil {
			return fmt.Sprintf("%v ", u.Username)
		}
	}
	return ""
}

func groupName(f os.FileInfo) string {
	if stat, ok := f.Sys().(*syscall.Stat_t); ok {
		if g, err := user.LookupGroupId(fmt.Sprint(stat.Gid)); err == nil {
			return fmt.Sprintf("%v ", g.Name)
		}
	}
	return ""
}

func linkCount(f os.FileInfo) string {
	if stat, ok := f.Sys().(*syscall.Stat_t); ok {
		return fmt.Sprintf("%v ", stat.Nlink)
	}
	return ""
}

func errCrossDevice(err error) bool {
	return err.(*os.LinkError).Err.(unix.Errno) == unix.EXDEV
}

func exportFiles(f string, fs []string, pwd string) {
	envFile := f
	envFiles := strings.Join(fs, genOpts.filesep)

	os.Setenv("f", envFile)
	os.Setenv("fs", envFiles)
	os.Setenv("PWD", pwd)

	if len(fs) == 0 {
		os.Setenv("fx", envFile)
	} else {
		os.Setenv("fx", envFiles)
	}
}
