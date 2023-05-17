package web

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/xnacly/manbib/database"
	"github.com/xnacly/manbib/shared"
)

//go:embed index.html
var index string
var tpl = template.Must(template.New("index.html").Parse(index))

func Run() {
	mux := http.NewServeMux()
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
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			// return
		}
		buf.WriteTo(w)
	})
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
	tpl.Execute(w, s)
}

func page(w http.ResponseWriter, r *http.Request) {

}
