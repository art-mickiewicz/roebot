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
var staging = make(map[int]State)

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

func setStateForID(id int, state State) {
	if state == Null || state == Clean { // Null and Clean states are not settable
		return
	}
	var oldState, newState State
	if checkState, ok := staging[id]; ok {
		oldState = checkState
	} else {
		oldState = Null
	}
	if oldState == newState {
		return
	}
	switch oldState {
	case Null:
		newState = state
	case Added:
		if state == Deleted {
			newState = Null
		} else {
			newState = oldState
		}
	case Updated, Deleted, Clean: // Pre-existente state
		if state == Added {
			newState = Updated
		} else {
			newState = state
		}
	}
	setStateForID(id, newState)
}

func SetTemplate(tpl Template) {
	var state State
	if _, ok := templates[tpl.ID]; ok {
		state = Updated
	} else {
		state = Added
	}
	templates[tpl.ID] = tpl
	staging[tpl.ID] = state
}

func DeleteTemplateByID(id int) int {
	was := len(templates)
	delete(templates, id)
	became := len(templates)
	if was-became > 0 {
		setStateForID(id, Deleted)
	}
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
