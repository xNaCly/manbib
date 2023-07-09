package main

import (
	"log"
	"os"

	"github.com/xnacly/manbib/cli"
	"github.com/xnacly/manbib/database"
	"github.com/xnacly/manbib/shared"
)

func main() {
	shared.Check()
	database.DB = database.Setup()
	defer database.DB.Conn.Close()

	if len(os.Args) == 1 {
		// TODO: indicate index status in frontend
		go func() {
			// TODO: pass --reindex to app for updating index
			if r, _ := database.DB.GetPagesAmount(); r != 0 {
				log.Println("skipping indexing, already indexed", r, "man pages")
				return
			}
			cli.Index()
		}()
		cli.StartWeb()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "cleardb":
		log.Println("clearing database...")
		database.DB.ClearDatabase()
	default:
		log.Fatalf("unknown instruction: '%s'", cmd)
	}
}
