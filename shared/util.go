package shared

import (
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
)

const HTML_TEMPLATE = `<nav id="$idprefix$TOC" role="doc-toc">
$if(toc-title)$
<h2 id="$idprefix$toc-title">$toc-title$</h2>
$endif$
$table-of-contents$
</nav>
$body$
`

var DEPENDENCIES = []string{
	"man",
	"zcat",
}

func Check() {
	if runtime.GOOS != "linux" {
		log.Fatalln("manbib only supports linux")
	}
	for _, d := range DEPENDENCIES {
		_, err := exec.LookPath(d)
		if err != nil {
			log.Fatalf("'%s' executable not found, please check its install.", d)
		}
	}
}

// returns the manbib config home, creates the folder if missing, creates pandoc template file
func ConfigHome() string {
	// ignoring this is safe, $HOME can't possibly be undefined
	v, _ := os.UserConfigDir()

	p := path.Join(v, "manbib")

	_, err := os.Stat(p)
	if err != nil {
		err = os.Mkdir(p, os.ModePerm)
		if err != nil {
			log.Fatalln("couldn't create manbib config home:", err)
		}
	}

	tPath := path.Join(p, "template.html5")
	_, err = os.Stat(tPath)

	if err != nil {
		err := os.WriteFile(tPath, []byte(HTML_TEMPLATE), os.ModePerm)
		if err != nil {
			log.Fatalln("couldn't create html5 pandoc template:", err)
		}
	}

	return p
}
