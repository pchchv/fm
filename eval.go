package main

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

func (e *setExpr) eval(app *app, args []string) {
	switch e.opt {
	case "anchorfind":
		genOpts.anchorfind = true
	case "noanchorfind":
		genOpts.anchorfind = false
	case "anchorfind!":
		genOpts.anchorfind = !genOpts.anchorfind
	case "autoquit":
		genOpts.autoquit = true
	case "noautoquit":
		genOpts.autoquit = false
	case "autoquit!":
		genOpts.autoquit = !genOpts.autoquit
	case "dircache":
		genOpts.dircache = true
	case "nodircache":
		genOpts.dircache = false
	case "dircache!":
		genOpts.dircache = !genOpts.dircache
	case "dircounts":
		genOpts.dircounts = true
	case "nodircounts":
		genOpts.dircounts = false
	case "dircounts!":
		genOpts.dircounts = !genOpts.dircounts
	case "dironly":
		genOpts.dironly = true
		app.nav.sort()
		app.nav.position()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "dirpreviews":
		genOpts.dirpreviews = true
	case "nodirpreviews":
		genOpts.dirpreviews = false
	case "dirpreviews!":
		genOpts.dirpreviews = !genOpts.dirpreviews
	case "nodironly":
		genOpts.dironly = false
		app.nav.sort()
		app.nav.position()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "dironly!":
		genOpts.dironly = !genOpts.dironly
		app.nav.sort()
		app.nav.position()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "dirfirst":
		genOpts.sortType.option |= dirfirstSort
		app.nav.sort()
		app.ui.sort()
	case "nodirfirst":
		genOpts.sortType.option &= ^dirfirstSort
		app.nav.sort()
		app.ui.sort()
	case "dirfirst!":
		genOpts.sortType.option ^= dirfirstSort
		app.nav.sort()
		app.ui.sort()
	case "drawbox":
		genOpts.drawbox = true
		app.ui.renew()
		if app.nav.height != app.ui.wins[0].h {
			app.nav.height = app.ui.wins[0].h
			app.nav.regCache = make(map[string]*reg)
		}
		app.ui.loadFile(app, true)
	case "nodrawbox":
		genOpts.drawbox = false
		app.ui.renew()
		if app.nav.height != app.ui.wins[0].h {
			app.nav.height = app.ui.wins[0].h
			app.nav.regCache = make(map[string]*reg)
		}
		app.ui.loadFile(app, true)
	case "drawbox!":
		genOpts.drawbox = !genOpts.drawbox
		app.ui.renew()
		if app.nav.height != app.ui.wins[0].h {
			app.nav.height = app.ui.wins[0].h
			app.nav.regCache = make(map[string]*reg)
		}
		app.ui.loadFile(app, true)
	case "globsearch":
		genOpts.globsearch = true
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "noglobsearch":
		genOpts.globsearch = false
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "globsearch!":
		genOpts.globsearch = !genOpts.globsearch
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "hidden":
		genOpts.sortType.option |= hiddenSort
		app.nav.sort()
		app.nav.position()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "nohidden":
		genOpts.sortType.option &= ^hiddenSort
		app.nav.sort()
		app.nav.position()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "hidden!":
		genOpts.sortType.option ^= hiddenSort
		app.nav.sort()
		app.nav.position()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "history":
		genOpts.history = true
	case "nohistory":
		genOpts.history = false
	case "history!":
		genOpts.history = !genOpts.history
	case "icons":
		genOpts.icons = true
	case "noicons":
		genOpts.icons = false
	case "icons!":
		genOpts.icons = !genOpts.icons
	case "ignorecase":
		genOpts.ignorecase = true
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "noignorecase":
		genOpts.ignorecase = false
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "ignorecase!":
		genOpts.ignorecase = !genOpts.ignorecase
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "ignoredia":
		genOpts.ignoredia = true
		app.nav.sort()
		app.ui.sort()
	case "noignoredia":
		genOpts.ignoredia = false
		app.nav.sort()
		app.ui.sort()
	case "ignoredia!":
		genOpts.ignoredia = !genOpts.ignoredia
		app.nav.sort()
		app.ui.sort()
	case "incfilter":
		genOpts.incfilter = true
	case "noincfilter":
		genOpts.incfilter = false
	case "incfilter!":
		genOpts.incfilter = !genOpts.incfilter
	case "incsearch":
		genOpts.incsearch = true
	case "noincsearch":
		genOpts.incsearch = false
	case "incsearch!":
		genOpts.incsearch = !genOpts.incsearch
	case "mouse":
		if !genOpts.mouse {
			genOpts.mouse = true
			app.ui.screen.EnableMouse(tcell.MouseButtonEvents)
		}
	case "nomouse":
		if genOpts.mouse {
			genOpts.mouse = false
			app.ui.screen.DisableMouse()
		}
	case "mouse!":
		if genOpts.mouse {
			genOpts.mouse = false
			app.ui.screen.DisableMouse()
		} else {
			genOpts.mouse = true
			app.ui.screen.EnableMouse(tcell.MouseButtonEvents)
		}
	case "number":
		genOpts.number = true
	case "nonumber":
		genOpts.number = false
	case "number!":
		genOpts.number = !genOpts.number
	case "preview":
		if len(genOpts.ratios) < 2 {
			app.ui.echoerr("preview: 'ratios' should consist of at least two numbers before enabling 'preview'")
			return
		}
		genOpts.preview = true
	case "nopreview":
		genOpts.preview = false
	case "preview!":
		if len(genOpts.ratios) < 2 {
			app.ui.echoerr("preview: 'ratios' should consist of at least two numbers before enabling 'preview'")
			return
		}
		genOpts.preview = !genOpts.preview
	case "relativenumber":
		genOpts.relativenumber = true
	case "norelativenumber":
		genOpts.relativenumber = false
	case "relativenumber!":
		genOpts.relativenumber = !genOpts.relativenumber
	case "reverse":
		genOpts.sortType.option |= reverseSort
		app.nav.sort()
		app.ui.sort()
	case "noreverse":
		genOpts.sortType.option &= ^reverseSort
		app.nav.sort()
		app.ui.sort()
	case "reverse!":
		genOpts.sortType.option ^= reverseSort
		app.nav.sort()
		app.ui.sort()
	case "smartcase":
		genOpts.smartcase = true
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "nosmartcase":
		genOpts.smartcase = false
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "smartcase!":
		genOpts.smartcase = !genOpts.smartcase
		app.nav.sort()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "smartdia":
		genOpts.smartdia = true
	case "nosmartdia":
		genOpts.smartdia = false
	case "smartdia!":
		genOpts.smartdia = !genOpts.smartdia
	case "waitmsg":
		genOpts.waitmsg = e.val
	case "wrapscan":
		genOpts.wrapscan = true
	case "nowrapscan":
		genOpts.wrapscan = false
	case "wrapscan!":
		genOpts.wrapscan = !genOpts.wrapscan
	case "wrapscroll":
		genOpts.wrapscroll = true
	case "nowrapscroll":
		genOpts.wrapscroll = false
	case "wrapscroll!":
		genOpts.wrapscroll = !genOpts.wrapscroll
	case "findlen":
		n, err := strconv.Atoi(e.val)
		if err != nil {
			app.ui.echoerrf("findlen: %s", err)
			return
		}
		if n < 0 {
			app.ui.echoerr("findlen: value should be a non-negative number")
			return
		}
		genOpts.findlen = n
	case "period":
		n, err := strconv.Atoi(e.val)
		if err != nil {
			app.ui.echoerrf("period: %s", err)
			return
		}
		if n < 0 {
			app.ui.echoerr("period: value should be a non-negative number")
			return
		}
		genOpts.period = n
		if n == 0 {
			app.ticker.Stop()
		} else {
			app.ticker.Stop()
			app.ticker = time.NewTicker(time.Duration(genOpts.period) * time.Second)
		}
	case "scrolloff":
		n, err := strconv.Atoi(e.val)
		if err != nil {
			app.ui.echoerrf("scrolloff: %s", err)
			return
		}
		if n < 0 {
			app.ui.echoerr("scrolloff: value should be a non-negative number")
			return
		}
		genOpts.scrolloff = n
	case "tabstop":
		n, err := strconv.Atoi(e.val)
		if err != nil {
			app.ui.echoerrf("tabstop: %s", err)
			return
		}
		if n <= 0 {
			app.ui.echoerr("tabstop: value should be a positive number")
			return
		}
		genOpts.tabstop = n
	case "errorfmt":
		genOpts.errorfmt = e.val
	case "filesep":
		genOpts.filesep = e.val
	case "hiddenfiles":
		toks := strings.Split(e.val, ":")
		for _, s := range toks {
			if s == "" {
				app.ui.echoerr("hiddenfiles: glob should be non-empty")
				return
			}
			_, err := filepath.Match(s, "a")
			if err != nil {
				app.ui.echoerrf("hiddenfiles: %s", err)
				return
			}
		}
		genOpts.hiddenfiles = toks
		app.nav.sort()
		app.nav.position()
		app.ui.sort()
		app.ui.loadFile(app, true)
	case "ifs":
		genOpts.ifs = e.val
	case "info":
		if e.val == "" {
			genOpts.info = nil
			return
		}
		toks := strings.Split(e.val, ":")
		for _, s := range toks {
			switch s {
			case "size", "time", "atime", "ctime":
			default:
				app.ui.echoerr("info: should consist of 'size', 'time', 'atime' or 'ctime' separated with colon")
				return
			}
		}
		genOpts.info = toks
	case "previewer":
		genOpts.previewer = replaceTilde(e.val)
	case "cleaner":
		genOpts.cleaner = replaceTilde(e.val)
	case "promptfmt":
		genOpts.promptfmt = e.val
	case "ratios":
		toks := strings.Split(e.val, ":")
		var rats []int
		for _, s := range toks {
			n, err := strconv.Atoi(s)
			if err != nil {
				app.ui.echoerrf("ratios: %s", err)
				return
			}
			if n <= 0 {
				app.ui.echoerr("ratios: value should be a positive number")
				return
			}
			rats = append(rats, n)
		}
		if genOpts.preview && len(rats) < 2 {
			app.ui.echoerr("ratios: should consist of at least two numbers when 'preview' is enabled")
			return
		}
		genOpts.ratios = rats
		app.ui.wins = getWins(app.ui.screen)
		app.ui.loadFile(app, true)
	case "selmode":
		genOpts.selmode = e.val
	case "shell":
		genOpts.shell = e.val
	case "shellflag":
		genOpts.shellflag = e.val
	case "shellopts":
		if e.val == "" {
			genOpts.shellopts = nil
			return
		}
		genOpts.shellopts = strings.Split(e.val, ":")
	case "sortby":
		switch e.val {
		case "natural":
			genOpts.sortType.method = naturalSort
		case "name":
			genOpts.sortType.method = nameSort
		case "size":
			genOpts.sortType.method = sizeSort
		case "time":
			genOpts.sortType.method = timeSort
		case "ctime":
			genOpts.sortType.method = ctimeSort
		case "atime":
			genOpts.sortType.method = atimeSort
		case "ext":
			genOpts.sortType.method = extSort
		default:
			app.ui.echoerr("sortby: value should either be 'natural', 'name', 'size', 'time', 'atime', 'ctime' or 'ext'")
			return
		}
		app.nav.sort()
		app.ui.sort()
	case "tempmarks":
		if e.val != "" {
			genOpts.tempmarks = "'" + e.val
		} else {
			genOpts.tempmarks = "'"
		}
	case "tagfmt":
		genOpts.tagfmt = e.val
	case "timefmt":
		genOpts.timefmt = e.val
	case "infotimefmtnew":
		genOpts.infotimefmtnew = e.val
	case "infotimefmtold":
		genOpts.infotimefmtold = e.val
	case "truncatechar":
		if runeSliceWidth([]rune(e.val)) != 1 {
			app.ui.echoerr("truncatechar: value should be a single character")
			return
		}

		genOpts.truncatechar = e.val
	default:
		// any key with the prefix user_ is accepted as a user defined option
		if strings.HasPrefix(e.opt, "user_") {
			genOpts.user[e.opt[5:]] = e.val
		} else {
			app.ui.echoerrf("unknown option: %s", e.opt)
		}
		return
	}
	app.ui.loadFileInfo(app.nav)
}

func (e *mapExpr) eval(app *app, args []string) {
	if e.expr == nil {
		delete(genOpts.keys, e.keys)
	} else {
		genOpts.keys[e.keys] = e.expr
	}
	app.ui.loadFileInfo(app.nav)
}

func (e *cmapExpr) eval(app *app, args []string) {
	if e.expr == nil {
		delete(genOpts.cmdkeys, e.key)
	} else {
		genOpts.cmdkeys[e.key] = e.expr
	}
	app.ui.loadFileInfo(app.nav)
}

func (e *cmdExpr) eval(app *app, args []string) {
	if e.expr == nil {
		delete(genOpts.cmds, e.name)
	} else {
		genOpts.cmds[e.name] = e.expr
	}
	app.ui.loadFileInfo(app.nav)
}

func preChdir(app *app) {
	if cmd, ok := genOpts.cmds["pre-cd"]; ok {
		cmd.eval(app, nil)
	}
}

func onChdir(app *app) {
	app.nav.addJumpList()
	if cmd, ok := genOpts.cmds["on-cd"]; ok {
		cmd.eval(app, nil)
	}
}

func onSelect(app *app) {
	if cmd, ok := genOpts.cmds["on-select"]; ok {
		cmd.eval(app, nil)
	}
}

func splitKeys(s string) (keys []string) {
	for i := 0; i < len(s); {
		r, w := utf8.DecodeRuneInString(s[i:])
		if r != '<' {
			keys = append(keys, s[i:i+w])
			i += w
		} else {
			j := i + w
			for r != '>' && j < len(s) {
				r, w = utf8.DecodeRuneInString(s[j:])
				j += w
			}
			keys = append(keys, s[i:j])
			i = j
		}
	}
	return
}

func doComplete(app *app) (matches []string) {
	switch app.ui.cmdPrefix {
	case ":":
		matches, app.ui.cmdAccLeft = completeCmd(app.ui.cmdAccLeft)
	case "/", "?":
		matches, app.ui.cmdAccLeft = completeFile(app.ui.cmdAccLeft)
	case "$", "%", "!", "&":
		matches, app.ui.cmdAccLeft = completeShell(app.ui.cmdAccLeft)
	}
	return
}

