package cli

import (
	"fmt"
	"log"
	"path"

	"github.com/xnacly/manbib/indexer"
	"github.com/xnacly/manbib/shared"
)

func Index() {
	log.Println("starting index creation, this may take a while")
	p := indexer.Lookup()
	l := len(p)
	log.Printf("found %d man pages, starting indexing", l)
	fmt.Println()
	templatePath := path.Join(shared.ConfigHome(), "template.html5")
	processedFiles := 0

	// this is extremly slow, split this up into around 9k areas and spawn a goroutine for each section
	for _, v := range p {
		fmt.Print("\033[1A\033[K")
		log.Printf("processed file: [%d/%d]\n", processedFiles, l)
		indexer.Index(v, templatePath)
		processedFiles++
	}
}
