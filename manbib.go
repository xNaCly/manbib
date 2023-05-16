package main

import (
	"log"
	"os"

	"github.com/xnacly/manbib/cli"
)

func main() {
	cli.Check()
	if len(os.Args) == 1 {
		cli.StartWeb()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "index", "i":
		cli.Index()
	default:
		log.Fatalf("unknown instruction: '%s'", cmd)
	}
}
