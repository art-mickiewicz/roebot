package state

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var schema = [...]string{
	"CREATE TABLE IF NOT EXISTS templates (id INTEGER PRIMARY KEY, source_message_id INTEGER NULL, target_message_id INTEGER NULL, text TEXT NOT NULL)",
}
var Templates = make([]Template, 10)
var Chats = make(map[string]Chat)

func init() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	for _, query := range schema {
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}
