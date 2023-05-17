package cli

import (
	"github.com/xnacly/manbib/web"
	"log"
)

func StartWeb() {
	log.Println("starting web interface: http://localhost:10997")
	web.Run()
}
