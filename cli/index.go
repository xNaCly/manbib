package cli

import (
	"fmt"
	"log"

	"github.com/xnacly/manbib/indexer"
)

func Index() {
	log.Println("starting indexing, this can take a while")
	p := indexer.Lookup()
	l := len(p)
	fmt.Println()
	for i, v := range p {
		fmt.Print("\033[1A\033[K")
		log.Printf("[%d/%d] starting thread for file. %s\n", i+1, l, v)
		// pass file to lexer, parser and code generator
	}
	// fill this into the database
}
