package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/xnacly/manbib/cli"
	"github.com/xnacly/manbib/roff"
)

func curTest() {
	r := roff.New()
	err := r.NewInput("test.1.gz")
	if err != nil {
		log.Fatalln(err)
	}
	v, _ := json.MarshalIndent(r.Start(), "", "  ")
	fmt.Println(string(v))
}

func main() {
	if len(os.Args) == 1 {
		cli.StartWeb()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "index", "i":
		cli.Index()
	case "dev":
		curTest()
	default:
		log.Fatalf("unknown instruction: '%s'", cmd)
	}
}
