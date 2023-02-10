package main

import (
	"os"
	"time"
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

func (file *file) TotalSize() int64 {
	if file.IsDir() {
		if file.dirSize >= 0 {
			return file.dirSize
		}
		return 0
	}
	return file.Size()
}
