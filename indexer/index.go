package indexer

import (
	"io/fs"
	"path/filepath"
)

var MAN_PATH = "/usr/share/man"

// queries the system for available man pages,
func Lookup() []string {
	r := make([]string, 0)
	filepath.WalkDir(MAN_PATH, func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() && len(path) != 0 {
			r = append(r, path)
		}
		return nil
	})
	return r
}
