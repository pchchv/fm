package main

import "time"

const (
	naturalSort sortMethod = iota
	nameSort
	sizeSort
	timeSort
	atimeSort
	ctimeSort
	extSort
)

const (
	dirfirstSort sortOption = 1 << iota
	hiddenSort
	reverseSort
)

var genOpts struct {
	anchorfind     bool
	autoquit       bool
	dircache       bool
	dircounts      bool
	dironly        bool
	dirpreviews    bool
	drawbox        bool
	globsearch     bool
	icons          bool
	ignorecase     bool
	ignoredia      bool
	incfilter      bool
	incsearch      bool
	mouse          bool
	number         bool
	preview        bool
	relativenumber bool
	smartcase      bool
	smartdia       bool
	waitmsg        string
	wrapscan       bool
	wrapscroll     bool
	findlen        int
	period         int
	scrolloff      int
	tabstop        int
	errorfmt       string
	filesep        string
	ifs            string
	previewer      string
	cleaner        string
	promptfmt      string
	selmode        string
	shell          string
	shellflag      string
	timefmt        string
	infotimefmtnew string
	infotimefmtold string
	truncatechar   string
	ratios         []int
	hiddenfiles    []string
	history        bool
	info           []string
	shellopts      []string
	keys           map[string]expr
	cmdkeys        map[string]expr
	cmds           map[string]expr
	user           map[string]string
	sortType       sortType
	tempmarks      string
	tagfmt         string
}

type sortMethod byte

type sortOption byte

type sortType struct {
	method sortMethod
	option sortOption
}

func init() {
	genOpts.anchorfind = true
	genOpts.autoquit = false
	genOpts.dircache = true
	genOpts.dircounts = false
	genOpts.dironly = false
	genOpts.dirpreviews = false
	genOpts.drawbox = false
	genOpts.globsearch = false
	genOpts.icons = false
	genOpts.ignorecase = true
	genOpts.ignoredia = true
	genOpts.incfilter = false
	genOpts.incsearch = false
	genOpts.mouse = false
	genOpts.number = false
	genOpts.preview = true
	genOpts.relativenumber = false
	genOpts.smartcase = true
	genOpts.smartdia = false
	genOpts.waitmsg = "Press any key to continue"
	genOpts.wrapscan = true
	genOpts.wrapscroll = false
	genOpts.findlen = 1
	genOpts.period = 0
	genOpts.scrolloff = 0
	genOpts.tabstop = 8
	genOpts.errorfmt = "\033[7;31;47m%s\033[0m"
	genOpts.filesep = "\n"
	genOpts.ifs = ""
	genOpts.previewer = ""
	genOpts.cleaner = ""
	genOpts.promptfmt = "\033[32;1m%u@%h\033[0m:\033[34;1m%d\033[0m\033[1m%f\033[0m"
	genOpts.selmode = "all"
	genOpts.shell = genDefaultShell
	genOpts.shellflag = genDefaultShellFlag
	genOpts.timefmt = time.ANSIC
	genOpts.infotimefmtnew = "Jan _2 15:04"
	genOpts.infotimefmtold = "Jan _2  2006"
	genOpts.truncatechar = "~"
	genOpts.ratios = []int{1, 2, 3}
	genOpts.hiddenfiles = []string{".*"}
	genOpts.history = true
	genOpts.info = nil
	genOpts.shellopts = nil
	genOpts.sortType = sortType{naturalSort, dirfirstSort}
	genOpts.tempmarks = "'"
	genOpts.tagfmt = "\033[31m%s\033[0m"

	genOpts.keys = make(map[string]expr)

	genOpts.keys["k"] = &callExpr{"up", nil, 1}
	genOpts.keys["<up>"] = &callExpr{"up", nil, 1}
	genOpts.keys["<m-up>"] = &callExpr{"up", nil, 1}
	genOpts.keys["<c-u>"] = &callExpr{"half-up", nil, 1}
	genOpts.keys["<c-b>"] = &callExpr{"page-up", nil, 1}
	genOpts.keys["<pgup>"] = &callExpr{"page-up", nil, 1}
	genOpts.keys["<c-y>"] = &callExpr{"scroll-up", nil, 1}
	genOpts.keys["j"] = &callExpr{"down", nil, 1}
	genOpts.keys["<down>"] = &callExpr{"down", nil, 1}
	genOpts.keys["<m-down>"] = &callExpr{"down", nil, 1}
	genOpts.keys["<c-d>"] = &callExpr{"half-down", nil, 1}
	genOpts.keys["<c-f>"] = &callExpr{"page-down", nil, 1}
	genOpts.keys["<pgdn>"] = &callExpr{"page-down", nil, 1}
	genOpts.keys["<c-e>"] = &callExpr{"scroll-down", nil, 1}
	genOpts.keys["h"] = &callExpr{"updir", nil, 1}
	genOpts.keys["<left>"] = &callExpr{"updir", nil, 1}
	genOpts.keys["l"] = &callExpr{"open", nil, 1}
	genOpts.keys["<right>"] = &callExpr{"open", nil, 1}
	genOpts.keys["q"] = &callExpr{"quit", nil, 1}
	genOpts.keys["gg"] = &callExpr{"top", nil, 1}
	genOpts.keys["<home>"] = &callExpr{"top", nil, 1}
	genOpts.keys["G"] = &callExpr{"bottom", nil, 1}
	genOpts.keys["<end>"] = &callExpr{"bottom", nil, 1}
	genOpts.keys["H"] = &callExpr{"high", nil, 1}
	genOpts.keys["M"] = &callExpr{"middle", nil, 1}
	genOpts.keys["L"] = &callExpr{"low", nil, 1}
	genOpts.keys["["] = &callExpr{"jump-prev", nil, 1}
	genOpts.keys["]"] = &callExpr{"jump-next", nil, 1}
	genOpts.keys["<space>"] = &listExpr{[]expr{&callExpr{"toggle", nil, 1}, &callExpr{"down", nil, 1}}, 1}
	genOpts.keys["t"] = &callExpr{"tag-toggle", nil, 1}
	genOpts.keys["v"] = &callExpr{"invert", nil, 1}
	genOpts.keys["u"] = &callExpr{"unselect", nil, 1}
	genOpts.keys["y"] = &callExpr{"copy", nil, 1}
	genOpts.keys["d"] = &callExpr{"cut", nil, 1}
	genOpts.keys["c"] = &callExpr{"clear", nil, 1}
	genOpts.keys["p"] = &callExpr{"paste", nil, 1}
	genOpts.keys["<c-l>"] = &callExpr{"redraw", nil, 1}
	genOpts.keys["<c-r>"] = &callExpr{"reload", nil, 1}
	genOpts.keys[":"] = &callExpr{"read", nil, 1}
	genOpts.keys["$"] = &callExpr{"shell", nil, 1}
	genOpts.keys["%"] = &callExpr{"shell-pipe", nil, 1}
	genOpts.keys["!"] = &callExpr{"shell-wait", nil, 1}
	genOpts.keys["&"] = &callExpr{"shell-async", nil, 1}
	genOpts.keys["f"] = &callExpr{"find", nil, 1}
	genOpts.keys["F"] = &callExpr{"find-back", nil, 1}
	genOpts.keys[";"] = &callExpr{"find-next", nil, 1}
	genOpts.keys[","] = &callExpr{"find-prev", nil, 1}
	genOpts.keys["/"] = &callExpr{"search", nil, 1}
	genOpts.keys["?"] = &callExpr{"search-back", nil, 1}
	genOpts.keys["n"] = &callExpr{"search-next", nil, 1}
	genOpts.keys["N"] = &callExpr{"search-prev", nil, 1}
	genOpts.keys["m"] = &callExpr{"mark-save", nil, 1}
	genOpts.keys["'"] = &callExpr{"mark-load", nil, 1}
	genOpts.keys[`"`] = &callExpr{"mark-remove", nil, 1}
	genOpts.keys[`r`] = &callExpr{"rename", nil, 1}
	genOpts.keys["<c-n>"] = &callExpr{"cmd-history-next", nil, 1}
	genOpts.keys["<c-p>"] = &callExpr{"cmd-history-prev", nil, 1}

	genOpts.keys["zh"] = &setExpr{"hidden!", ""}
	genOpts.keys["zr"] = &setExpr{"reverse!", ""}
	genOpts.keys["zn"] = &setExpr{"info", ""}
	genOpts.keys["zs"] = &setExpr{"info", "size"}
	genOpts.keys["zt"] = &setExpr{"info", "time"}
	genOpts.keys["za"] = &setExpr{"info", "size:time"}
	genOpts.keys["sn"] = &listExpr{[]expr{&setExpr{"sortby", "natural"}, &setExpr{"info", ""}}, 1}
	genOpts.keys["ss"] = &listExpr{[]expr{&setExpr{"sortby", "size"}, &setExpr{"info", "size"}}, 1}
	genOpts.keys["st"] = &listExpr{[]expr{&setExpr{"sortby", "time"}, &setExpr{"info", "time"}}, 1}
	genOpts.keys["sa"] = &listExpr{[]expr{&setExpr{"sortby", "atime"}, &setExpr{"info", "atime"}}, 1}
	genOpts.keys["sc"] = &listExpr{[]expr{&setExpr{"sortby", "ctime"}, &setExpr{"info", "ctime"}}, 1}
	genOpts.keys["se"] = &listExpr{[]expr{&setExpr{"sortby", "ext"}, &setExpr{"info", ""}}, 1}
	genOpts.keys["gh"] = &callExpr{"cd", []string{"~"}, 1}

	genOpts.cmdkeys = make(map[string]expr)

	genOpts.cmdkeys["<space>"] = &callExpr{"cmd-insert", []string{" "}, 1}
	genOpts.cmdkeys["<esc>"] = &callExpr{"cmd-escape", nil, 1}
	genOpts.cmdkeys["<tab>"] = &callExpr{"cmd-complete", nil, 1}
	genOpts.cmdkeys["<enter>"] = &callExpr{"cmd-enter", nil, 1}
	genOpts.cmdkeys["<c-j>"] = &callExpr{"cmd-enter", nil, 1}
	genOpts.cmdkeys["<c-n>"] = &callExpr{"cmd-history-next", nil, 1}
	genOpts.cmdkeys["<c-p>"] = &callExpr{"cmd-history-prev", nil, 1}
	genOpts.cmdkeys["<delete>"] = &callExpr{"cmd-delete", nil, 1}
	genOpts.cmdkeys["<c-d>"] = &callExpr{"cmd-delete", nil, 1}
	genOpts.cmdkeys["<backspace>"] = &callExpr{"cmd-delete-back", nil, 1}
	genOpts.cmdkeys["<backspace2>"] = &callExpr{"cmd-delete-back", nil, 1}
	genOpts.cmdkeys["<left>"] = &callExpr{"cmd-left", nil, 1}
	genOpts.cmdkeys["<c-b>"] = &callExpr{"cmd-left", nil, 1}
	genOpts.cmdkeys["<right>"] = &callExpr{"cmd-right", nil, 1}
	genOpts.cmdkeys["<c-f>"] = &callExpr{"cmd-right", nil, 1}
	genOpts.cmdkeys["<home>"] = &callExpr{"cmd-home", nil, 1}
	genOpts.cmdkeys["<c-a>"] = &callExpr{"cmd-home", nil, 1}
	genOpts.cmdkeys["<end>"] = &callExpr{"cmd-end", nil, 1}
	genOpts.cmdkeys["<c-e>"] = &callExpr{"cmd-end", nil, 1}
	genOpts.cmdkeys["<c-u>"] = &callExpr{"cmd-delete-home", nil, 1}
	genOpts.cmdkeys["<c-k>"] = &callExpr{"cmd-delete-end", nil, 1}
	genOpts.cmdkeys["<c-w>"] = &callExpr{"cmd-delete-unix-word", nil, 1}
	genOpts.cmdkeys["<c-y>"] = &callExpr{"cmd-yank", nil, 1}
	genOpts.cmdkeys["<c-t>"] = &callExpr{"cmd-transpose", nil, 1}
	genOpts.cmdkeys["<c-c>"] = &callExpr{"cmd-interrupt", nil, 1}
	genOpts.cmdkeys["<a-f>"] = &callExpr{"cmd-word", nil, 1}
	genOpts.cmdkeys["<a-b>"] = &callExpr{"cmd-word-back", nil, 1}
	genOpts.cmdkeys["<a-c>"] = &callExpr{"cmd-capitalize-word", nil, 1}
	genOpts.cmdkeys["<a-d>"] = &callExpr{"cmd-delete-word", nil, 1}
	genOpts.cmdkeys["<a-u>"] = &callExpr{"cmd-uppercase-word", nil, 1}
	genOpts.cmdkeys["<a-l>"] = &callExpr{"cmd-lowercase-word", nil, 1}
	genOpts.cmdkeys["<a-t>"] = &callExpr{"cmd-transpose-word", nil, 1}

	genOpts.cmds = make(map[string]expr)
	genOpts.user = make(map[string]string)

	setDefaults()
}
