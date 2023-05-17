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

// TODO: add this to the welcome page
type Stats struct {
	Total int
	Man1  int
	Man2  int
	Man3  int
	Man4  int
	Man5  int
	Man6  int
	Man7  int
	Man8  int
}

type SearchTemplateContent struct {
	Query        string
	Rows         []Page
	ResultAmount int
	Latency      string
	Total        int // TODO: see shared/models.go:L14
	Page         Page
}
