package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djherbis/times"
	"github.com/pchchv/golog"
)

const (
	notLink linkState = iota
	working
	broken
)

type linkState byte

type file struct {
	os.FileInfo
	linkState  linkState
	linkTarget string
	path       string
	dirCount   int
	dirSize    int64
	accessTime time.Time
	changeTime time.Time
	ext        string
}

type dir struct {
	loading     bool      // directory is loading from disk
	loadTime    time.Time // current loading or last load time
	ind         int       // index of current entry in files
	pos         int       // position of current entry in ui
	path        string    // full path of directory
	files       []*file   // displayed files in directory including or excluding hidden ones
	allFiles    []*file   // all files in directory including hidden ones (same array as files)
	sortType    sortType  // sort method and options from last sort
	dironly     bool      // dironly value from last sort
	hiddenfiles []string  // hiddenfiles value from last sort
	filter      []string  // last filter for this directory
	ignorecase  bool      // ignorecase value from last sort
	ignoredia   bool      // ignoredia value from last sort
	noPerm      bool      // whether lf has no permission to open the directory
	lines       []string  // lines of text to display if directory previews are enabled
}

type indexedSelections struct {
	paths   []string
	indices []int
}

func newDir(path string) *dir {
	time := time.Now()

	files, err := readdir(path)
	if err != nil {
		golog.Info("reading directory: %s", err)
	}

	return &dir{
		loading:  genOpts.dirpreviews, // directory is loaded after previewer function exits.
		loadTime: time,
		path:     path,
		files:    files,
		allFiles: files,
		noPerm:   os.IsPermission(err),
	}
}

func (file *file) TotalSize() int64 {
	if file.IsDir() {
		if file.dirSize >= 0 {
			return file.dirSize
		}
		return 0
	}
	return file.Size()
}

func (dir *dir) name() string {
	if len(dir.files) == 0 {
		return ""
	}

	return dir.files[dir.ind].Name()
}

func (dir *dir) sel(name string, height int) {
	if len(dir.files) == 0 {
		dir.ind, dir.pos = 0, 0
		return
	}

	dir.ind = max(dir.ind, 0)
	dir.ind = min(dir.ind, len(dir.files)-1)

	if dir.files[dir.ind].Name() != name {
		for i, f := range dir.files {
			if f.Name() == name {
				dir.ind = i
				break
			}
		}
	}

	edge := min(min(height/2, genOpts.scrolloff), len(dir.files)-dir.ind-1)
	dir.pos = min(dir.ind, height-edge-1)
}

func (dir *dir) sort() {
	dir.sortType = genOpts.sortType
	dir.dironly = genOpts.dironly
	dir.hiddenfiles = genOpts.hiddenfiles
	dir.ignorecase = genOpts.ignorecase
	dir.ignoredia = genOpts.ignoredia
	dir.files = dir.allFiles

	switch dir.sortType.method {
	case naturalSort:
		sort.SliceStable(dir.files, func(i, j int) bool {
			s1, s2 := normalize(dir.files[i].Name(), dir.files[j].Name(), dir.ignorecase, dir.ignoredia)
			return naturalLess(s1, s2)
		})
	case nameSort:
		sort.SliceStable(dir.files, func(i, j int) bool {
			s1, s2 := normalize(dir.files[i].Name(), dir.files[j].Name(), dir.ignorecase, dir.ignoredia)
			return s1 < s2
		})
	case sizeSort:
		sort.SliceStable(dir.files, func(i, j int) bool {
			return dir.files[i].TotalSize() < dir.files[j].TotalSize()
		})
	case timeSort:
		sort.SliceStable(dir.files, func(i, j int) bool {
			return dir.files[i].ModTime().Before(dir.files[j].ModTime())
		})
	case atimeSort:
		sort.SliceStable(dir.files, func(i, j int) bool {
			return dir.files[i].accessTime.Before(dir.files[j].accessTime)
		})
	case ctimeSort:
		sort.SliceStable(dir.files, func(i, j int) bool {
			return dir.files[i].changeTime.Before(dir.files[j].changeTime)
		})
	case extSort:
		sort.SliceStable(dir.files, func(i, j int) bool {
			ext1, ext2 := normalize(dir.files[i].ext, dir.files[j].ext, dir.ignorecase, dir.ignoredia)

			// if the extension cannot be defined (directories, files without an extension),
			// a null byte is used so that these files can be ranked higher
			if ext1 == "" {
				ext1 = "\x00"
			}
			if ext2 == "" {
				ext2 = "\x00"
			}

			name1, name2 := normalize(dir.files[i].Name(), dir.files[j].Name(), dir.ignorecase, dir.ignoredia)

			// for natural sorting of filenames, the name is combined with ext, but ext must come first
			return ext1 < ext2 || ext1 == ext2 && name1 < name2
		})
	}

	if dir.sortType.option&reverseSort != 0 {
		for i, j := 0, len(dir.files)-1; i < j; i, j = i+1, j-1 {
			dir.files[i], dir.files[j] = dir.files[j], dir.files[i]
		}
	}

	if dir.sortType.option&dirfirstSort != 0 {
		sort.SliceStable(dir.files, func(i, j int) bool {
			if dir.files[i].IsDir() == dir.files[j].IsDir() {
				return i < j
			}
			return dir.files[i].IsDir()
		})
	}

	// when the dironly option is enabled, the files are moved to the beginning of our file list,
	// and then the beginning of the displayed files is set to the first directory in the list
	if dir.dironly {
		sort.SliceStable(dir.files, func(i, j int) bool {
			if !dir.files[i].IsDir() && !dir.files[j].IsDir() {
				return i < j
			}
			return !dir.files[i].IsDir()
		})
		dir.files = func() []*file {
			for i, f := range dir.files {
				if f.IsDir() {
					return dir.files[i:]
				}
			}
			return dir.files[len(dir.files):]
		}()
	}

	// when the hide option is disabled, hidden files are moved to the beginning of the file list,
	// and then the beginning of the displayed files is set to the first unhidden file in the list
	if dir.sortType.option&hiddenSort == 0 {
		sort.SliceStable(dir.files, func(i, j int) bool {
			if isHidden(dir.files[i], dir.path, dir.hiddenfiles) && isHidden(dir.files[j], dir.path, dir.hiddenfiles) {
				return i < j
			}
			return isHidden(dir.files[i], dir.path, dir.hiddenfiles)
		})
		for i, f := range dir.files {
			if !isHidden(f, dir.path, dir.hiddenfiles) {
				dir.files = dir.files[i:]
				break
			}
		}
		if len(dir.files) > 0 && isHidden(dir.files[len(dir.files)-1], dir.path, dir.hiddenfiles) {
			dir.files = dir.files[len(dir.files):]
		}
	}

	if len(dir.filter) != 0 {
		sort.SliceStable(dir.files, func(i, j int) bool {
			if isFiltered(dir.files[i], dir.filter) && isFiltered(dir.files[j], dir.filter) {
				return i < j
			}
			return isFiltered(dir.files[i], dir.filter)
		})
		for i, f := range dir.files {
			if !isFiltered(f, dir.filter) {
				dir.files = dir.files[i:]
				break
			}
		}
		if len(dir.files) > 0 && isFiltered(dir.files[len(dir.files)-1], dir.filter) {
			dir.files = dir.files[len(dir.files):]
		}
	}

	dir.ind = max(dir.ind, 0)
	dir.ind = min(dir.ind, len(dir.files)-1)
}
func (m indexedSelections) Len() int {
	return len(m.paths)
}

func (m indexedSelections) Swap(i, j int) {
	m.paths[i], m.paths[j] = m.paths[j], m.paths[i]
	m.indices[i], m.indices[j] = m.indices[j], m.indices[i]
}

func (m indexedSelections) Less(i, j int) bool {
	return m.indices[i] < m.indices[j]
}

func readdir(path string) ([]*file, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()

	files := make([]*file, 0, len(names))
	for _, fname := range names {
		var linkState linkState
		var linkTarget string

		fpath := filepath.Join(path, fname)

		lstat, err := os.Lstat(fpath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			golog.Info("getting file information: %s", err)
			continue
		}

		if lstat.Mode()&os.ModeSymlink != 0 {
			stat, err := os.Stat(fpath)
			if err == nil {
				linkState = working
				lstat = stat
			} else {
				linkState = broken
			}
			linkTarget, err = os.Readlink(fpath)
			if err != nil {
				golog.Info("reading link target: %s", err)
			}
		}

		ts := times.Get(lstat)
		at := ts.AccessTime()
		var ct time.Time
		// from times docs: ChangeTime() panics unless HasChangeTime() is true
		if ts.HasChangeTime() {
			ct = ts.ChangeTime()
		} else {
			// fall back to ModTime if ChangeTime cannot be determined
			ct = lstat.ModTime()
		}

		// returns an empty string if extension could not be determined
		// i.e. directories, filenames without extensions
		ext := filepath.Ext(fpath)

		dirCount := -1
		if lstat.IsDir() && genOpts.dircounts {
			d, err := os.Open(fpath)
			if err != nil {
				dirCount = -2
			} else {
				names, err := d.Readdirnames(1000)
				d.Close()

				if names == nil && err != io.EOF {
					dirCount = -2
				} else {
					dirCount = len(names)
				}
			}
		}

		files = append(files, &file{
			FileInfo:   lstat,
			linkState:  linkState,
			linkTarget: linkTarget,
			path:       fpath,
			dirCount:   dirCount,
			dirSize:    -1,
			accessTime: at,
			changeTime: ct,
			ext:        ext,
		})
	}

	return files, err
}

func normalize(s1, s2 string, ignorecase, ignoredia bool) (string, string) {
	if genOpts.ignorecase {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}
	if genOpts.ignoredia {
		s1 = removeDiacritics(s1)
		s2 = removeDiacritics(s2)
	}
	return s1, s2
}

func searchMatch(name, pattern string) (matched bool, err error) {
	if genOpts.ignorecase {
		lpattern := strings.ToLower(pattern)
		if !genOpts.smartcase || lpattern == pattern {
			pattern = lpattern
			name = strings.ToLower(name)
		}
	}
	if genOpts.ignoredia {
		lpattern := removeDiacritics(pattern)
		if !genOpts.smartdia || lpattern == pattern {
			pattern = lpattern
			name = removeDiacritics(name)
		}
	}
	if genOpts.globsearch {
		return filepath.Match(pattern, name)
	}
	return strings.Contains(name, pattern), nil
}

//lint:ignore U1000 This function is not used on Windows
func matchPattern(pattern, name, path string) bool {
	s := name

	pattern = replaceTilde(pattern)

	if filepath.IsAbs(pattern) {
		s = filepath.Join(path, name)
	}

	// pattern errors are checked when 'hiddenfiles' option is set
	matched, _ := filepath.Match(pattern, s)

	return matched
}

func findMatch(name, pattern string) bool {
	if genOpts.ignorecase {
		lpattern := strings.ToLower(pattern)
		if !genOpts.smartcase || lpattern == pattern {
			pattern = lpattern
			name = strings.ToLower(name)
		}
	}
	if genOpts.ignoredia {
		lpattern := removeDiacritics(pattern)
		if !genOpts.smartdia || lpattern == pattern {
			pattern = lpattern
			name = removeDiacritics(name)
		}
	}
	if genOpts.anchorfind {
		return strings.HasPrefix(name, pattern)
	}
	return strings.Contains(name, pattern)
}

func isFiltered(f os.FileInfo, filter []string) bool {
	for _, pattern := range filter {
		matched, err := searchMatch(f.Name(), strings.TrimPrefix(pattern, "!"))
		if err != nil {
			log.Printf("Filter Error: %s", err)
			return false
		}
		if strings.HasPrefix(pattern, "!") && matched {
			return true
		} else if !strings.HasPrefix(pattern, "!") && !matched {
			return true
		}
	}
	return false
}
