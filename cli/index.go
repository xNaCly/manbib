package cli

import (
	"fmt"
	"log"
	"path"
	"time"

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

	// TODO: this is extremly slow, split this up into around 9k areas and spawn a goroutine for each section
	start := time.Now()
	for i, v := range p {
		fmt.Print("\033[1A\033[K")
		log.Printf("processed file: [%d/%d] %.2f %%\n", i, l, (float64(i)/float64(l))*100.0)
		indexer.Index(v, templatePath)
	}
	log.Println("done, took: ", time.Since(start))
}
