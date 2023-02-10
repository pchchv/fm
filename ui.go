package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/pchchv/golog"
)

const (
	genEscapeCode = 27

	// These colors are not currently customizeable
	LineNumberColor = tcell.ColorOlive
	SelectionColor  = tcell.ColorPurple
	YankColor       = tcell.ColorOlive
	CutColor        = tcell.ColorMaroon
)

var genThisYear = time.Now().Year()

type win struct {
	w, h, x, y int
}

type dirContext struct {
	selections map[string]int
	saves      map[string]bool
	tags       map[string]string
}

type dirStyle struct {
	colors     styleMap
	icons      iconMap
	previewing bool
}

type reg struct {
	loading  bool
	volatile bool
	loadTime time.Time
	path     string
	lines    []string
}

func newWin(w, h, x, y int) *win {
	return &win{w, h, x, y}
}

func (win *win) renew(w, h, x, y int) {
	win.w, win.h, win.x, win.y = w, h, x, y
}

func (win *win) print(screen tcell.Screen, x, y int, st tcell.Style, s string) tcell.Style {
	off := x
	var comb []rune

	for i := 0; i < len(s); i++ {
		r, w := utf8.DecodeRuneInString(s[i:])

		if r == genEscapeCode && i+1 < len(s) && s[i+1] == '[' {
			j := strings.IndexAny(s[i:min(len(s), i+64)], "mK")
			if j == -1 {
				continue
			}
			if s[i+j] == 'm' {
				st = applyAnsiCodes(s[i+2:i+j], st)
			}

			i += j
			continue
		}

		for {
			rc, wc := utf8.DecodeRuneInString(s[i+w:])
			if !unicode.Is(unicode.Mn, rc) {
				break
			}
			comb = append(comb, rc)
			i += wc
		}

		if x < win.w {
			screen.SetContent(win.x+x, win.y+y, r, comb, st)
			comb = nil
		}

		i += w - 1

		if r == '\t' {
			s := genOpts.tabstop - (x-off)%genOpts.tabstop
			for i := 0; i < s && x+i < win.w; i++ {
				screen.SetContent(win.x+x+i, win.y+y, ' ', nil, st)
			}
			x += s
		} else {
			x += runewidth.RuneWidth(r)
		}
	}

	return st
}

func (win *win) printf(screen tcell.Screen, x, y int, st tcell.Style, format string, a ...interface{}) {
	win.print(screen, x, y, st, fmt.Sprintf(format, a...))
}

func (win *win) printLine(screen tcell.Screen, x, y int, st tcell.Style, s string) {
	win.printf(screen, x, y, st, "%s%*s", s, win.w-printLength(s), "")
}

func (win *win) printRight(screen tcell.Screen, y int, st tcell.Style, s string) {
	win.print(screen, win.w-printLength(s), y, st, s)
}

func (win *win) printReg(screen tcell.Screen, reg *reg) {
	if reg == nil {
		return
	}

	st := tcell.StyleDefault

	if reg.loading {
		st = st.Reverse(true)
		win.print(screen, 2, 0, st, "loading...")
		return
	}

	for i, l := range reg.lines {
		if i > win.h-1 {
			break
		}

		st = win.print(screen, 2, i, st, l)
	}
}

func (win *win) printDir(screen tcell.Screen, dir *dir, context *dirContext, dirStyle *dirStyle) {
	var lnwidth int
	var lnformat string

	if win.w < 5 || dir == nil {
		return
	}

	messageStyle := tcell.StyleDefault.Reverse(true)

	if dir.noPerm {
		win.print(screen, 2, 0, messageStyle, "permission denied")
		return
	}
	if (dir.loading && len(dir.files) == 0) || (dirStyle.previewing && dir.loading && genOpts.dirpreviews) {
		win.print(screen, 2, 0, messageStyle, "loading...")
		return
	}

	if dirStyle.previewing && genOpts.dirpreviews && len(genOpts.previewer) > 0 {
		// Print previewer result instead of default directory print operation.
		st := tcell.StyleDefault
		for i, l := range dir.lines {
			if i > win.h-1 {
				break
			}

			st = win.print(screen, 2, i, st, l)
		}
		return
	}
	if len(dir.files) == 0 {
		win.print(screen, 2, 0, messageStyle, "empty")
		return
	}

	beg := max(dir.ind-dir.pos, 0)
	end := min(beg+win.h, len(dir.files))

	if beg > end {
		return
	}

	if genOpts.number || genOpts.relativenumber {
		lnwidth = 1
		if genOpts.number && genOpts.relativenumber {
			lnwidth++
		}
		for j := 10; j < len(dir.files); j *= 10 {
			lnwidth++
		}
		lnformat = fmt.Sprintf("%%%d.d ", lnwidth)
	}

	for i, f := range dir.files[beg:end] {
		st := dirStyle.colors.get(f)

		if lnwidth > 0 {
			var ln string

			if genOpts.number && (!genOpts.relativenumber) {
				ln = fmt.Sprintf(lnformat, i+1+beg)
			} else if genOpts.relativenumber {
				switch {
				case i < dir.pos:
					ln = fmt.Sprintf(lnformat, dir.pos-i)
				case i > dir.pos:
					ln = fmt.Sprintf(lnformat, i-dir.pos)
				case genOpts.number:
					ln = fmt.Sprintf(fmt.Sprintf("%%%d.d ", lnwidth-1), i+1+beg)
				default:
					ln = ""
				}
			}

			win.print(screen, 0, i, tcell.StyleDefault.Foreground(LineNumberColor), ln)
		}

		path := filepath.Join(dir.path, f.Name())

		if _, ok := context.selections[path]; ok {
			win.print(screen, lnwidth, i, st.Background(SelectionColor), " ")
		} else if cp, ok := context.saves[path]; ok {
			if cp {
				win.print(screen, lnwidth, i, st.Background(YankColor), " ")
			} else {
				win.print(screen, lnwidth, i, st.Background(CutColor), " ")
			}
		}

		if i == dir.pos {
			st = st.Reverse(true)
		}

		var s []rune

		s = append(s, ' ')

		var iwidth int

		if genOpts.icons {
			s = append(s, []rune(dirStyle.icons.get(f))...)
			s = append(s, ' ')
			iwidth = 2
		}

		for _, r := range f.Name() {
			s = append(s, r)
		}

		w := runeSliceWidth(s)

		if w > win.w-3 {
			s = runeSliceWidthRange(s, 0, win.w-4)
			s = append(s, []rune(genOpts.truncatechar)...)
		} else {
			for i := 0; i < win.w-3-w; i++ {
				s = append(s, ' ')
			}
		}

		info := fileInfo(f, dir)

		if len(info) > 0 && win.w-lnwidth-iwidth-2 > 2*len(info) {
			if win.w-2 > w+len(info) {
				s = runeSliceWidthRange(s, 0, win.w-3-len(info)-lnwidth)
			} else {
				s = runeSliceWidthRange(s, 0, win.w-4-len(info)-lnwidth)
				s = append(s, []rune(genOpts.truncatechar)...)
			}
			for _, r := range info {
				s = append(s, r)
			}
		}

		s = append(s, ' ')

		win.print(screen, lnwidth+1, i, st, string(s))

		tag, ok := context.tags[path]
		if ok {
			if i == dir.pos {
				win.print(screen, lnwidth+1, i, st, tag)
			} else {
				win.print(screen, lnwidth+1, i, tcell.StyleDefault, fmt.Sprintf(genOpts.tagfmt, tag))
			}
		}
	}
}

func printLength(s string) int {
	ind := 0
	off := 0
	for i := 0; i < len(s); i++ {
		r, w := utf8.DecodeRuneInString(s[i:])

		if r == genEscapeCode && i+1 < len(s) && s[i+1] == '[' {
			j := strings.IndexAny(s[i:min(len(s), i+64)], "mK")
			if j == -1 {
				continue
			}

			i += j
			continue
		}

		i += w - 1

		if r == '\t' {
			ind += genOpts.tabstop - (ind-off)%genOpts.tabstop
		} else {
			ind += runewidth.RuneWidth(r)
		}
	}

	return ind
}

func fileInfo(f *file, d *dir) string {
	var info string

	for _, s := range genOpts.info {
		switch s {
		case "size":
			if !(f.IsDir() && genOpts.dircounts) {
				var sz string
				if f.IsDir() && f.dirSize < 0 {
					sz = "-"
				} else {
					sz = humanize(f.TotalSize())
				}
				info = fmt.Sprintf("%s %4s", info, sz)
				continue
			}

			switch {
			case f.dirCount < -1:
				info = fmt.Sprintf("%s    !", info)
			case f.dirCount < 0:
				info = fmt.Sprintf("%s    ?", info)
			case f.dirCount < 1000:
				info = fmt.Sprintf("%s %4d", info, f.dirCount)
			default:
				info = fmt.Sprintf("%s 999+", info)
			}
		case "time":
			info = fmt.Sprintf("%s %*s", info, max(len(genOpts.infotimefmtnew), len(genOpts.infotimefmtold)), infotimefmt(f.ModTime()))
		case "atime":
			info = fmt.Sprintf("%s %*s", info, max(len(genOpts.infotimefmtnew), len(genOpts.infotimefmtold)), infotimefmt(f.accessTime))
		case "ctime":
			info = fmt.Sprintf("%s %*s", info, max(len(genOpts.infotimefmtnew), len(genOpts.infotimefmtold)), infotimefmt(f.changeTime))
		default:
			golog.Info("unknown info type: %s", s)
		}
	}

	return info
}

func infotimefmt(t time.Time) string {
	if t.Year() == genThisYear {
		return t.Format(genOpts.infotimefmtnew)
	}
	return t.Format(genOpts.infotimefmtold)
}
