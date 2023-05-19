package web

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/xnacly/manbib/database"
	"github.com/xnacly/manbib/shared"
)

//go:embed style.css
var style []byte

//go:embed index.html
var indexTemplate string
var tpl = template.Must(template.New("index.html").Parse(indexTemplate))

//go:embed page.html
var pageTemplate string
var pageTpl = template.Must(template.New("page.html").Parse(pageTemplate))

//go:embed search.html
var searchTemplate string
var searchTpl = template.Must(template.New("search.html").Parse(searchTemplate))

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write(style)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		randomPage, err := database.DB.GetRandomPage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		indexedAmount, err := database.DB.GetPagesAmount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, shared.SearchTemplateContent{
			Page:  randomPage,
			Total: indexedAmount,
		})
		// INFO: we ignore errors while executing the template,
		// this keeps everything neat and simple
		// https://open.spotify.com/track/3koCCeSaVUyrRo3N2gHrd8
		// if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
		// }
		buf.WriteTo(w)
	})
	mux.HandleFunc("/page", page)
	mux.HandleFunc("/search", search)
	log.Fatalln(http.ListenAndServe(":10997", mux))
}

func search(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := url.Query()
	query := params.Get("q")
	limit := params.Get("l")

	qStr := ""
	if len(query) != 0 {
		qStr = fmt.Sprintf("%s%%", query)
	}

	if len(limit) == 0 {
		limit = "200"
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()
	rows := database.DB.GetPages(qStr, l)
	s := shared.SearchTemplateContent{
		Query:        query,
		Rows:         rows,
		ResultAmount: len(rows),
		Latency:      time.Since(start).String(),
	}
	buf := &bytes.Buffer{}
	err = searchTpl.Execute(buf, s)
	buf.WriteTo(w)
}

// TODO: do this on the fly, if a user opens the manpage, convert it to html
// TODO: override the preview in the database, not only here
func page(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := url.Query()
	reqPage := params.Get("p")
	page := database.DB.GetPages(reqPage, 1)
	templatePath := path.Join(shared.ConfigHome(), "template.html5")
	if len(page) < 1 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	p := page[0]
	if len(p.Preview) == 0 {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("zcat %s | pandoc --from man --to markdown | pandoc --toc --from markdown --to html5 --template %s", p.Path, templatePath))
		manPreview, err := cmd.Output()
		if err != nil {
			log.Println(err, string(manPreview), cmd.String())
		}
		p.Preview = string(manPreview)
	}

	buf := &bytes.Buffer{}
	err = pageTpl.Execute(buf, p)
	buf.WriteTo(w)
}
