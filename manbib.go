package main

import (
	"flag"
	"log"

	"github.com/xnacly/manbib/cli"
	"github.com/xnacly/manbib/database"
	"github.com/xnacly/manbib/shared"
	"github.com/xnacly/manbib/web"
)

func main() {
	reindex := flag.Bool("reindex", false, "recreate the man page index")
	port := flag.Int("port", 10997, "port to start web interface on")
	flag.Parse()

	shared.Check()
	database.DB = database.Setup()
	defer database.DB.Conn.Close()

	if *reindex {
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

	log.Printf("starting web interface: http://localhost:%d", *port)
	web.Run(*port)
	return
}
