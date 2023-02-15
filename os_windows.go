package main

import "golang.org/x/sys/windows"

var (
	envOpener = os.Getenv("OPENER")
	envEditor = os.Getenv("EDITOR")
	envPager  = os.Getenv("PAGER")
	envShell  = os.Getenv("SHELL")

	envPathExt = os.Getenv("PATHEXT")

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

func init() {
	if envOpener == "" {
		envOpener = `start ""`
	}

	if envEditor == "" {
		envEditor = "notepad"
	}

	if envPager == "" {
		envPager = "more"
	}

	if envShell == "" {
		envShell = "cmd"
	}

	u, err := user.Current()
	if err != nil {
		log.Printf("user: %s", err)
	}
	genUser = u

	// remove domain prefix
	genUser.Username = strings.Split(genUser.Username, `\`)[1]

	data := os.Getenv("LOCALAPPDATA")

	genConfigPaths = []string{
		filepath.Join(os.Getenv("ProgramData"), "fm", "fmrc"),
		filepath.Join(data, "fm", "fmrc"),
	}

	genColorsPaths = []string{
		filepath.Join(os.Getenv("ProgramData"), "fm", "colors"),
		filepath.Join(data, "fm", "colors"),
	}

	genIconsPaths = []string{
		filepath.Join(os.Getenv("ProgramData"), "fm", "icons"),
		filepath.Join(data, "fm", "icons"),
	}

	genFilesPath = filepath.Join(data, "fm", "files")
	genMarksPath = filepath.Join(data, "fm", "marks")
	genTagsPath = filepath.Join(data, "fm", "tags")
	genHistoryPath = filepath.Join(data, "fm", "history")
}

func detachedCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &windows.SysProcAttr{CreationFlags: 8}
	return cmd
}

func shellCommand(s string, args []string) *exec.Cmd {
	args = append([]string{gOpts.shellflag, s}, args...)

	args = append(gOpts.shellopts, args...)

	return exec.Command(gOpts.shell, args...)
}

func shellSetPG(cmd *exec.Cmd) {
}

func shellKill(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

func setDefaults() {
	gOpts.cmds["open"] = &execExpr{"&", "%OPENER% %f%"}
	gOpts.keys["e"] = &execExpr{"$", "%EDITOR% %f%"}
	gOpts.keys["i"] = &execExpr{"!", "%PAGER% %f%"}
	gOpts.keys["w"] = &execExpr{"$", "%SHELL%"}

	gOpts.cmds["doc"] = &execExpr{"!", "fm -doc | %PAGER%"}
	gOpts.keys["<f-1>"] = &callExpr{"doc", nil, 1}
}

func setUserUmask() {}

func isExecutable(f os.FileInfo) bool {
	exts := strings.Split(envPathExt, string(filepath.ListSeparator))
	for _, e := range exts {
		if strings.HasSuffix(strings.ToLower(f.Name()), strings.ToLower(e)) {
			log.Print(f.Name(), e)
			return true
		}
	}
	return false
}

func isHidden(f os.FileInfo, path string, hiddenfiles []string) bool {
	ptr, err := windows.UTF16PtrFromString(filepath.Join(path, f.Name()))
	if err != nil {
		return false
	}
	attrs, err := windows.GetFileAttributes(ptr)
	if err != nil {
		return false
	}
	return attrs&windows.FILE_ATTRIBUTE_HIDDEN != 0
}

func userName(f os.FileInfo) string {
	return ""
}

func groupName(f os.FileInfo) string {
	return ""
}

func linkCount(f os.FileInfo) string {
	return ""
}

func errCrossDevice(err error) bool {
	return err.(*os.LinkError).Err.(windows.Errno) == 17
}

func exportFiles(f string, fs []string, pwd string) {
	envFile := fmt.Sprintf(`"%s"`, f)

	var quotedFiles []string
	for _, f := range fs {
		quotedFiles = append(quotedFiles, fmt.Sprintf(`"%s"`, f))
	}
	envFiles := strings.Join(quotedFiles, gOpts.filesep)

	envPWD := fmt.Sprintf(`"%s"`, pwd)

	os.Setenv("f", envFile)
	os.Setenv("fs", envFiles)
	os.Setenv("PWD", envPWD)

	if len(fs) == 0 {
		os.Setenv("fx", envFile)
	} else {
		os.Setenv("fx", envFiles)
	}
}
