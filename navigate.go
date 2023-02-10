package main

import (
	"io"
	"os"
	"path/filepath"
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
