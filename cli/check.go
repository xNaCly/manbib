package cli

import (
	"log"
	"os/exec"
	"runtime"
)

func Check() {
	if runtime.GOOS != "linux" {
		log.Fatalln("manbib only supports linux")
	}
	_, err := exec.LookPath("pandoc")
	if err != nil {
		log.Fatalln("'pandoc' executable not found, please check your installation.")
	}
}
