package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type cmdItem struct {
	prefix string
	value  string
}

type app struct {
	ui             *ui
	nav            *nav
	ticker         *time.Ticker
	quitChan       chan struct{}
	cmd            *exec.Cmd
	cmdIn          io.WriteCloser
	cmdOutBuf      []byte
	cmdHistory     []cmdItem
	cmdHistoryBeg  int
	cmdHistoryInd  int
	menuCompActive bool
	menuComps      []string
	menuCompInd    int
}

func newApp(ui *ui, nav *nav) *app {
	quitChan := make(chan struct{}, 1)

	app := &app{
		ui:       ui,
		nav:      nav,
		ticker:   new(time.Ticker),
		quitChan: quitChan,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		switch <-sigChan {
		case os.Interrupt:
			return
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM:
			app.quit()
			os.Exit(3)
			return
		}
	}()

	return app
}

func (app *app) quit() {
	if genOpts.history {
		if err := app.writeHistory(); err != nil {
			log.Printf("writing history file: %s", err)
		}
	}
	if !genSingleMode {
		if err := remote(fmt.Sprintf("drop %d", genClientID)); err != nil {
			log.Printf("dropping connection: %s", err)
		}
		if genOpts.autoquit {
			if err := remote("quit"); err != nil {
				log.Printf("auto quitting server: %s", err)
			}
		}
	}
}

func loadFiles() (list []string, cp bool, err error) {
	files, err := os.Open(genFilesPath)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		err = fmt.Errorf("opening file selections file: %s", err)
		return
	}
	defer files.Close()

	s := bufio.NewScanner(files)

	s.Scan()

	switch s.Text() {
	case "copy":
		cp = true
	case "move":
		cp = false
	default:
		err = fmt.Errorf("unexpected option to copy file(s): %s", s.Text())
		return
	}

	for s.Scan() && s.Text() != "" {
		list = append(list, s.Text())
	}

	if s.Err() != nil {
		err = fmt.Errorf("scanning file list: %s", s.Err())
		return
	}

	log.Printf("loading files: %v", list)

	return
}

func (app *app) readHistory() error {
	f, err := os.Open(genHistoryPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("opening history file: %s", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		toks := strings.SplitN(scanner.Text(), " ", 2)
		if toks[0] != ":" && toks[0] != "$" && toks[0] != "%" && toks[0] != "!" && toks[0] != "&" {
			continue
		}
		if len(toks) < 2 {
			continue
		}
		app.cmdHistory = append(app.cmdHistory, cmdItem{toks[0], toks[1]})
	}

	app.cmdHistoryBeg = len(app.cmdHistory)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading history file: %s", err)
	}

	return nil
}

func (app *app) writeHistory() error {
	if len(app.cmdHistory) == 0 {
		return nil
	}

	local := make([]cmdItem, len(app.cmdHistory)-app.cmdHistoryBeg)
	copy(local, app.cmdHistory[app.cmdHistoryBeg:])
	app.cmdHistory = nil

	if err := app.readHistory(); err != nil {
		return fmt.Errorf("reading history file: %s", err)
	}

	app.cmdHistory = append(app.cmdHistory, local...)

	if err := os.MkdirAll(filepath.Dir(genHistoryPath), os.ModePerm); err != nil {
		return fmt.Errorf("creating data directory: %s", err)
	}

	f, err := os.Create(genHistoryPath)
	if err != nil {
		return fmt.Errorf("creating history file: %s", err)
	}
	defer f.Close()

	if len(app.cmdHistory) > 1000 {
		app.cmdHistory = app.cmdHistory[len(app.cmdHistory)-1000:]
	}

	for _, cmd := range app.cmdHistory {
		_, err = f.WriteString(fmt.Sprintf("%s %s\n", cmd.prefix, cmd.value))
		if err != nil {
			return fmt.Errorf("writing history file: %s", err)
		}
	}

	return nil
}
