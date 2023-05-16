// holds a tiny sqlite3 wrapper written using the database/sql package with the go-sqlite3 lib,
// contains methods which allow for easy interaction with the database
//
// Usage Example:
//
//	database.DB = database.Setup()
//
//	database.DB.InsertPage(shared.Page{
//	    Path:        "/usr/share/man/man1/ls.1.gz",
//	    Name:        "ls",
//	    Preview:     "lool ls",
//	    LastUpdated: time.Now(),
//	})
//
//	r, _ := json.MarshalIndent(database.DB.GetPages("%ls%"), "", "\t")
//	fmt.Println(string(r))
//	database.DB.ClearDatabase()
//	r, _ = json.MarshalIndent(database.DB.GetPages(""), "", "\t")
//	fmt.Println(string(r))
package database

// TODO: add shared.Page.Examples & shared.Page.History

import (
	"database/sql"
	"fmt"
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

	return Database{
		Conn: db,
	}
}

// drops the pages table
func (d *Database) ClearDatabase() {
	d.Conn.Exec("DROP TABLE pages")
}

// inserts the given page in the pages table of the database
func (d *Database) InsertPage(p shared.Page) {
	_, err := d.Conn.Exec("INSERT INTO pages (path, name, preview, last_updated) VALUES(?, ?, ?, ?)", p.Path, p.Name, p.Preview, p.LastUpdated)
	if err != nil {
		return
	}
}

// queries the database with the given name via a WHERE name LIKE  statement
func (d *Database) GetPages(name string) []shared.Page {
	var rows *sql.Rows
	var err error
	if len(name) != 0 {
		rows, err = d.Conn.Query("SELECT path, name, preview, last_updated FROM pages WHERE name LIKE ?", name)
	} else {
		rows, err = d.Conn.Query("SELECT path, name, preview, last_updated FROM pages")
	}
	res := make([]shared.Page, 0)
	if err != nil {
		return res
	}
	for rows.Next() {
		r := shared.Page{}
		err = rows.Scan(&r.Path, &r.Name, &r.Preview, &r.LastUpdated)
		if err != nil {
			fmt.Println("hit:", err)
			continue
		}
		res = append(res, r)
	}
	return res
}
