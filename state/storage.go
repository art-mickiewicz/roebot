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

//var templates = make([]Template, 0, 10)
var templates = make(map[int]Template)

func NewTemplate(targetChannel string, srcPtr MessagePtr, text string) Template {
	maxID := 0
	for _, t := range templates {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newID := maxID + 1
	return Template{ID: newID, TargetChannel: targetChannel, SourceMessagePtr: srcPtr, Text: text}
}

func GetTemplateBySource(srcPtr MessagePtr) (tpl Template, ok bool) {
	for k, t := range templates {
		if t.SourceMessagePtr == srcPtr {
			return templates[k], true
		}
	}
	return Template{}, false
}

func GetTemplateByID(id int) (tpl Template, ok bool) {
	for i, t := range templates {
		if t.ID == id {
			return templates[i], true
		}
	}
	return Template{}, false
}

func GetTemplatesCount() int {
	return len(templates)
}

func GetTemplates() []Template {
	tpls := make([]Template, len(templates))
	i := 0
	for _, v := range templates {
		tpls[i] = v
		i++
	}
	return tpls
}

func SetTemplate(tpl Template) {
	templates[tpl.ID] = tpl
}

func DeleteTemplateByID(id int) int {
	was := len(templates)
	delete(templates, id)
	became := len(templates)
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
