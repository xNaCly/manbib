// holds a tiny sqlite3 wrapper written using the database/sql package with the go-sqlite3 lib,
// contains methods which allow for easy interaction with the database
//
// Usage Example:
//
//	database.DB = database.Setup()
//
//	r, _ := json.MarshalIndent(database.DB.GetPages("%ls%"), "", "\t")
//	fmt.Println(string(r))
//	database.DB.ClearDatabase()
//	r, _ = json.MarshalIndent(database.DB.GetPages(""), "", "\t")
//	fmt.Println(string(r))
package database

// TODO: add shared.Page.Examples & shared.Page.History
// Examples from tldr, History from default shell ($SHELL)

import (
	"database/sql"
	"log"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xnacly/manbib/shared"
)

type Database struct {
	Conn *sql.DB // underlying database connection
}

var DB Database

// opens db connection, tries ping, creates pages table, errors in either of the three operations result in exiting the program
func Setup() Database {
	db, err := sql.Open("sqlite3", path.Join(shared.ConfigHome(), "manbib.db"))

	if err != nil {
		log.Fatalln("failed to establish database connection:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("couldn't reach database:", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS pages (id INTEGER PRIMARY KEY, path text UNIQUE, name text, preview text, examples text, history text, last_updated timestamp)")

	if err != nil {
		log.Fatalln("failed to create 'pages' table:", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY, page_id INTEGER, searched_at timestamp)")

	if err != nil {
		log.Fatalln("failed to create 'history' table:", err)
	}

	return Database{
		Conn: db,
	}
}

// drops the pages table
func (d *Database) ClearDatabase() {
	d.Conn.Exec("DROP TABLE pages")
	d.Conn.Exec("DROP TABLE history")
}

// updates the preview of the page with the given path
func (d *Database) UpdatePreview(path string, preview []byte) error {
	_, err := d.Conn.Exec("UPDATE pages SET preview = ? WHERE path = ?", preview, path)
	return err
}

// inserts the given page in the pages table of the database
func (d *Database) InsertPages(p []shared.Page) error {
	t, err := d.Conn.Begin()
	if err != nil {
		return err
	}

	for _, page := range p {
		_, err := t.Exec("INSERT INTO pages (path, name, preview, last_updated) VALUES(?, ?, ?, ?)", page.Path, page.Name, page.Preview, page.LastUpdated)
		if err != nil {
			return err
		}
	}

	err = t.Commit()
	return err
}

// queries the database with the given name via a WHERE name LIKE  statement
func (d *Database) GetPages(name string, limit int) []shared.Page {
	var rows *sql.Rows
	var err error
	if limit < 0 {
		limit = 1
	}
	if len(name) != 0 {
		rows, err = d.Conn.Query("SELECT path, name, preview, last_updated FROM pages WHERE name LIKE ? LIMIT ?", name, limit)
	} else {
		rows, err = d.Conn.Query("SELECT path, name, preview, last_updated FROM pages LIMIT ?", limit)
	}
	res := make([]shared.Page, 0)
	if err != nil {
		log.Println(err)
		return res
	}
	for rows.Next() {
		r := shared.Page{}
		err = rows.Scan(&r.Path, &r.Name, &r.Preview, &r.LastUpdated)
		if err != nil {
			continue
		}
		res = append(res, r)
	}
	return res
}

func (d *Database) GetPagesAmount() (int, error) {
	row := d.Conn.QueryRow("SELECT COUNT(*) as p FROM pages")
	var p int
	err := row.Scan(&p)
	if err != nil {
		return 0, err
	}
	return p, nil
}

func (d *Database) GetRandomPage() (shared.Page, error) {
	row := d.Conn.QueryRow("SELECT path, name, preview, last_updated FROM pages ORDER BY RANDOM() LIMIT 1")
	r := shared.Page{}

	// FIXES: errors if the result set is empty
	if row == nil {
		return shared.Page{}, nil
	}

	err := row.Scan(&r.Path, &r.Name, &r.Preview, &r.LastUpdated)
	if err != nil {
		return shared.Page{}, err
	}
	return r, nil
}
