package state

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var schema = [...]string{
	`CREATE TABLE IF NOT EXISTS templates
		( id INTEGER PRIMARY KEY
		, source_chat_id INTEGER NULL
		, source_message_id INTEGER NULL
		, target_chat_id INTEGER NULL
		, target_message_id INTEGER NULL
		, text TEXT NOT NULL
		)`,
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	for _, query := range schema {
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}

func PersistTemplates() {
	/* Deleted */
	deletedTemplates := GetTemplatesWithState(Deleted)
	if len(deletedTemplates) > 0 {
		deletedIds := make([]string, len(deletedTemplates))
		for i, tpl := range deletedTemplates {
			deletedIds[i] = string(tpl.ID)
			delete(staging, tpl.ID)
		}
		inExpr := strings.Join(deletedIds, ",")
		if _, err := db.Exec("DELETE FROM templates WHERE id IN (" + inExpr + ")"); err != nil {
			log.Fatal(err)
		}
	}

	/* Updated */
	updatedTemplates := GetTemplatesWithState(Updated)
	if len(updatedTemplates) > 0 {
		stmt, err := db.Prepare(`
			UPDATE templates
			SET source_chat_id=?, source_message_id=?, target_chat_id=?, target_message_id=?, text=?
			WHERE id=?
		`)
		if err != nil {
			log.Fatal(err)
		}
		for _, tpl := range updatedTemplates {
			_, err := stmt.Exec(
				tpl.SourceMessagePtr.ChatID, tpl.SourceMessagePtr.MessageID,
				tpl.TargetMessagePtr.ChatID, tpl.TargetMessagePtr.MessageID,
				tpl.Text, tpl.ID,
			)
			if err != nil {
				log.Fatal(err)
			}
			setStateForID(tpl.ID, Clean)
		}
	}

	/* Added */
	addedTemplates := GetTemplatesWithState(Deleted)
	if len(addedTemplates) > 0 {
		stmt, err := db.Prepare(`
			INSERT INTO templates
			(id, source_chat_id, source_message_id, target_chat_id, target_message_id, text)
			VALUES (?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			log.Fatal(err)
		}
		for _, tpl := range updatedTemplates {
			_, err := stmt.Exec(
				tpl.ID,
				tpl.SourceMessagePtr.ChatID, tpl.SourceMessagePtr.MessageID,
				tpl.TargetMessagePtr.ChatID, tpl.TargetMessagePtr.MessageID,
				tpl.Text,
			)
			if err != nil {
				log.Fatal(err)
			}
			setStateForID(tpl.ID, Clean)
		}
	}
}

func LoadTemplates() {
	rows, err := db.Query(`
		SELECT id, source_chat_id, source_message_id, target_chat_id, target_message_id, text
		FROM templates
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		tpl := Template{}
		err := rows.Scan(
			&(tpl.ID), &(tpl.SourceMessagePtr.ChatID), &(tpl.SourceMessagePtr.MessageID),
			&(tpl.TargetMessagePtr.ChatID), &(tpl.TargetMessagePtr.MessageID), &(tpl.Text),
		)
		if err != nil {
			log.Fatal(err)
		}
		templates[tpl.ID] = tpl
	}
}
