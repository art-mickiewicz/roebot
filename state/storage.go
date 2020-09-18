package state

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"robpike.io/filter"
)

var db *sql.DB
var schema = [...]string{
	"CREATE TABLE IF NOT EXISTS templates (id INTEGER PRIMARY KEY, source_message_id INTEGER NULL, target_message_id INTEGER NULL, text TEXT NOT NULL)",
}
var Templates = make([]Template, 0, 10)

func NewTemplate(targetChannel string, srcPtr MessagePtr, text string) Template {
	maxID := 0
	for _, t := range Templates {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newID := maxID + 1
	return Template{ID: newID, TargetChannel: targetChannel, SourceMessagePtr: srcPtr, Text: text}
}

func GetTemplateBySource(srcPtr MessagePtr) *Template {
	for i, t := range Templates {
		if t.SourceMessagePtr == srcPtr {
			return &Templates[i]
		}
	}
	return nil
}

func GetTemplateByID(id int) *Template {
	for i, t := range Templates {
		if t.ID == id {
			return &Templates[i]
		}
	}
	return nil
}

func DeleteTemplateByID(id int) int {
	was := len(Templates)
	Templates = filter.Choose(Templates, func(t Template) bool {
		return t.ID != id
	}).([]Template)
	became := len(Templates)
	return was - became
}

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
