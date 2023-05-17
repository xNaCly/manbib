package shared

import "time"

type Page struct {
	Path        string
	Name        string
	Preview     string
	Examples    string
	History     string
	LastUpdated time.Time
}

type SearchTemplateContent struct {
	Query        string
	Rows         []Page
	ResultAmount int
	Latency      string
	Total        int
	Page         Page
}
