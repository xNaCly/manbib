package indexer

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/xnacly/manbib/database"
	"github.com/xnacly/manbib/shared"
)

var MAN_PATH = "/usr/share/man/"

// queries the system for available man pages,
func Lookup() (r []string) {
	r = make([]string, 0)
	filepath.WalkDir(MAN_PATH, func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() {
			r = append(r, path)
		}
		return nil
	})
	return
}

func Index(p string, templatePath string) {
	// TODO: do this on the fly, if a user opens the manpage, convert it to html
	// pandocCmd := fmt.Sprintf("man -Thtml %s", p)
	// manPreview, _ := exec.Command("bash", "-c", pandocCmd).Output()

	cmdName := strings.Replace(path.Base(p), ".gz", "", 1)
	database.DB.InsertPage(shared.Page{
		Name:        cmdName,
		Path:        p,
		Preview:     "",
		LastUpdated: time.Now(),
	})
}
