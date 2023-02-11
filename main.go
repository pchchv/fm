package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/pchchv/golog"
)

var (
	envPath  = os.Getenv("PATH")
	envLevel = os.Getenv("FM_LEVEL")

	genSingleMode    bool
	genClientID      int
	genHostname      string
	genLastDirPath   string
	genSelectionPath string
	genSocketProt    string
	genSocketPath    string
	genLogPath       string
	genSelect        string
	genConfigPath    string
	genCommands      arrayFlag
	genVersion       string
	genDocString     = `TODO`
)

type arrayFlag []string

func (a *arrayFlag) Set(v string) error {
	*a = append(*a, v)
	return nil
}

func (a *arrayFlag) String() string {
	return strings.Join(*a, ", ")
}

func startServer() {
	cmd := detachedCommand(os.Args[0], "-server")
	if err := cmd.Start(); err != nil {
		golog.Info("starting server: %s", err)
	}
}

func checkServer() {
	if genSocketProt == "unix" {
		if _, err := os.Stat(genSocketPath); os.IsNotExist(err) {
			startServer()
		} else if _, err := net.Dial(genSocketProt, genSocketPath); err != nil {
			os.Remove(genSocketPath)
			startServer()
		}
	} else {
		if _, err := net.Dial(genSocketProt, genSocketPath); err != nil {
			startServer()
		}
	}
}

func exportEnvVars() {
	os.Setenv("id", strconv.Itoa(genClientID))

	os.Setenv("OPENER", envOpener)
	os.Setenv("EDITOR", envEditor)
	os.Setenv("PAGER", envPager)
	os.Setenv("SHELL", envShell)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getting current directory: %s\n", err)
	}
	os.Setenv("OLDPWD", dir)

	level, err := strconv.Atoi(envLevel)
	if err != nil {
		log.Printf("reading lf level: %s", err)
	}

	level++

	os.Setenv("FM_LEVEL", strconv.Itoa(level))
}

func fieldToString(field reflect.Value) string {
	kind := field.Kind()
	var value string

	switch kind {
	case reflect.Int:
		value = strconv.Itoa(int(field.Int()))
	case reflect.Bool:
		value = strconv.FormatBool(field.Bool())
	case reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			element := field.Index(i)

			if i == 0 {
				value = fieldToString(element)
			} else {
				value += ":" + fieldToString(element)
			}
		}
	default:
		value = field.String()
	}

	return value
}

func exportOpts() {
	e := reflect.ValueOf(&genOpts).Elem()

	for i := 0; i < e.NumField(); i++ {
		// Get name and prefix it with lf_
		name := e.Type().Field(i).Name
		name = fmt.Sprintf("lf_%s", name)

		// Skip maps
		if name == "lf_keys" || name == "lf_cmdkeys" || name == "lf_cmds" {
			continue
		}

		// Get string representation of the value
		if name == "lf_sortType" {
			var sortby string

			switch genOpts.sortType.method {
			case naturalSort:
				sortby = "natural"
			case nameSort:
				sortby = "name"
			case sizeSort:
				sortby = "size"
			case timeSort:
				sortby = "time"
			case ctimeSort:
				sortby = "ctime"
			case atimeSort:
				sortby = "atime"
			case extSort:
				sortby = "ext"
			}

			os.Setenv("lf_sortby", sortby)

			reverse := strconv.FormatBool(genOpts.sortType.option&reverseSort != 0)
			os.Setenv("lf_reverse", reverse)

			hidden := strconv.FormatBool(genOpts.sortType.option&hiddenSort != 0)
			os.Setenv("lf_hidden", hidden)

			dirfirst := strconv.FormatBool(genOpts.sortType.option&dirfirstSort != 0)
			os.Setenv("lf_dirfirst", dirfirst)
		} else if name == "lf_user" {
			// set each user option
			for key, value := range genOpts.user {
				os.Setenv(name+"_"+key, value)
			}
		} else {
			field := e.Field(i)
			value := fieldToString(field)

			os.Setenv(name, value)
		}
	}
}

func main() {
	flag.Usage = func() {
		f := flag.CommandLine.Output()
		fmt.Fprintln(f, "lf - Terminal file manager")
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

	flag.StringVar(&genLastDirPath, "last-dir-path", "", "path to the file to write the last dir on exit (to use for cd)")
	flag.StringVar(&genSelectionPath, "selection-path", "", "path to the file to write selected files on open (to use as open file dialog)")
	flag.StringVar(&genConfigPath, "config", "", "path to the config file (instead of the usual paths)")
	flag.Var(&genCommands, "command", "command to execute on client initialization")
	flag.StringVar(&genLogPath, "log", "", "path to the log file to write messages")

	flag.Parse()

	genSocketProt = genDefaultSocketProt
	genSocketPath = genDefaultSocketPath

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			golog.Fatal("could not create CPU profile: %s", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			golog.Fatal("could not start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
	}

	switch {
	case *showDoc:
		fmt.Println(genDocString)
	case *showVersion:
		fmt.Println(genVersion)
	case *remoteCmd != "":
		if err := remote(*remoteCmd); err != nil {
			golog.Fatal("remote command: %s", err)
		}
	case *serverMode:
		if genLogPath != "" && !filepath.IsAbs(genLogPath) {
			wd, err := os.Getwd()
			if err != nil {
				golog.Fatal("getting current directory: %s", err)
			} else {
				genLogPath = filepath.Join(wd, genLogPath)
			}
		}
		os.Chdir(genUser.HomeDir)
		serve()
	default:
		genSingleMode = *singleMode

		if !genSingleMode {
			checkServer()
		}

		genClientID = os.Getpid()

		switch flag.NArg() {
		case 0:
			_, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "getting current directory: %s\n", err)
				os.Exit(2)
			}
		case 1:
			genSelect = flag.Arg(0)
		default:
			fmt.Fprintf(os.Stderr, "only single file or directory is allowed\n")
			os.Exit(2)
		}

		exportEnvVars()

		run()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			golog.Fatal("could not create memory profile: ", err)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			golog.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
