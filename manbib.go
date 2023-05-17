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
		cli.StartWeb()
		return
	}

	// TODO: remove index cmd, index when user starts webfrontend
	cmd := os.Args[1]
	switch cmd {
	case "index", "i":
		cli.Index()
	case "cleardb":
		log.Println("clearing database...")
		database.DB.ClearDatabase()
	default:
		log.Fatalf("unknown instruction: '%s'", cmd)
	}
}
