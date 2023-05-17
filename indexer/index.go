package indexer

import (
	"fmt"
	"io/fs"
	"os/exec"
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
	// pandocCmd := fmt.Sprintf("zcat %s | pandoc --from man --to markdown | pandoc --toc --from markdown --to html5 --template %s", p, templatePath)
	pandocCmd := fmt.Sprintf("zcat %s | pandoc --from man --to html", p)
	manPreview, _ := exec.Command("bash", "-c", pandocCmd).Output()
	cmdName := strings.Replace(path.Base(p), ".gz", "", 1)
	database.DB.InsertPage(shared.Page{
		Name:        cmdName,
		Path:        p,
		Preview:     string(manPreview),
		LastUpdated: time.Now(),
	})
}
