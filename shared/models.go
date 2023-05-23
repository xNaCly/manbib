package shared

import "time"

// represents a man page in the database
type Page struct {
	Id          int       // identification
	Path        string    // source file path: /usr/share/man/man1/ls.1.gz
	Name        string    // name of the executable: ls.1
	Preview     string    // html preview of the man page: <h1>ls.1</h1>
	Examples    string    // usage examples
	History     string    // past usages
	LastUpdated time.Time // timestamp of last indexing
}

// represents an opened man page of the past in the database
type HistoryItem struct {
	PageId     int
	SearchedAt time.Time
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
