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
