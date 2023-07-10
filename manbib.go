package main

import (
	"flag"
	"log"

	"github.com/xnacly/manbib/cli"
	"github.com/xnacly/manbib/database"
	"github.com/xnacly/manbib/shared"
)

func main() {
	REINDEX := flag.Bool("reindex", false, "specify to recreate the index")
	flag.Parse()

	shared.Check()
	database.DB = database.Setup()
	defer database.DB.Conn.Close()

	if *REINDEX {
		r, _ := database.DB.GetPagesAmount()
		log.Printf("removing '%d' pages", r)
		database.DB.ClearDatabase()
		log.Println("recreating index")
	} else {
		// TODO: indicate index status in frontend
		go func() {
			if r, _ := database.DB.GetPagesAmount(); r != 0 {
				log.Println("skipping indexing, already indexed", r, "man pages")
				return
			}
			cli.Index()
		}()
	}

	cli.StartWeb()
	return
}
