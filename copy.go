package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func copySize(srcs []string) (int64, error) {
	var total int64

	for _, src := range srcs {
		_, err := os.Lstat(src)
		if os.IsNotExist(err) {
			return total, fmt.Errorf("src does not exist: %q", src)
		}

		err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk: %s", err)
			}
			total += info.Size()
			return nil
		})

		if err != nil {
			return total, err
		}
	}

	return total, nil
}
