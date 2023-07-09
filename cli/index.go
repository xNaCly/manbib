package cli

import (
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/xnacly/manbib/database"
	"github.com/xnacly/manbib/indexer"
	"github.com/xnacly/manbib/shared"
)

func Index() {
	log.Println("starting index creation, this may take a while")
	p := indexer.Lookup()
	l := len(p)
	r := make([]shared.Page, 0)
	log.Printf("found %d man pages, starting indexing", l)
	fmt.Println()

	start := time.Now()
	for i, v := range p {
		fmt.Print("\033[1A\033[K")
		i += 1
		log.Printf("processed file: [%d/%d] %.2f %%\n", i, l, (float64(i)/float64(l))*100.0)

		r = append(r, shared.Page{
			Name:        strings.Replace(path.Base(v), ".gz", "", 1),
			Path:        v,
			Preview:     "",
			LastUpdated: start,
		})
	}
	err := database.DB.InsertPages(r)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("done, took: ", time.Since(start))
}
