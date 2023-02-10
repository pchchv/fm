package main

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
