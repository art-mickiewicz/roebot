package state

import (
	"database/sql"
	"log"
	"strconv"
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
		, channel_name TEXT NULL
		, text TEXT NOT NULL
		)`,
	`CREATE TABLE IF NOT EXISTS chats
		( id INTEGER PRIMARY KEY
		, username TEXT NOT NULL
		, title TEXT NOT NULL
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
			deletedIds[i] = strconv.Itoa(tpl.ID)
			setStateForID(tpl.ID, Null)
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
			SET source_chat_id=?, source_message_id=?, target_chat_id=?, target_message_id=?, channel_name=?, text=?
			WHERE id=?
		`)
		if err != nil {
			log.Fatal(err)
		}
		for _, tpl := range updatedTemplates {
			_, err := stmt.Exec(
				tpl.SourceMessagePtr.ChatID, tpl.SourceMessagePtr.MessageID,
				tpl.TargetMessagePtr.ChatID, tpl.TargetMessagePtr.MessageID,
				tpl.TargetChannel, tpl.Text, tpl.ID,
			)
			if err != nil {
				log.Fatal(err)
			}
			setStateForID(tpl.ID, Clean)
		}
	}

	/* Added */
	addedTemplates := GetTemplatesWithState(Added)
	if len(addedTemplates) > 0 {
		stmt, err := db.Prepare(`
			INSERT INTO templates
			(id, source_chat_id, source_message_id, target_chat_id, target_message_id, channel_name, text)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			log.Fatal(err)
		}
		for _, tpl := range addedTemplates {
			_, err := stmt.Exec(
				tpl.ID,
				tpl.SourceMessagePtr.ChatID, tpl.SourceMessagePtr.MessageID,
				tpl.TargetMessagePtr.ChatID, tpl.TargetMessagePtr.MessageID,
				tpl.TargetChannel, tpl.Text,
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
		SELECT id, source_chat_id, source_message_id, target_chat_id, target_message_id, channel_name, text
		FROM templates ORDER BY id
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		tpl := Template{}
		err := rows.Scan(
			&(tpl.ID), &(tpl.SourceMessagePtr.ChatID), &(tpl.SourceMessagePtr.MessageID),
			&(tpl.TargetMessagePtr.ChatID), &(tpl.TargetMessagePtr.MessageID), &(tpl.TargetChannel),
			&(tpl.Text),
		)
		if err != nil {
			log.Fatal(err)
		}
		tpl, err = parseTemplate(tpl)
		if err == nil {
			templates[tpl.ID] = tpl
		}
	}
}

func PersistChat(chat Chat, update bool) {
	if update {
		_, err := db.Exec("UPDATE chats SET username = ?, title = ? WHERE id = ?",
			chat.Username, chat.Title, chat.ID)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.Exec("INSERT INTO chats (id, username, title) VALUES (?, ?, ?)",
			chat.ID, chat.Username, chat.Title)
		if err != nil {
			log.Println(err)
		}
	}
}

func LoadChats() {
	rows, err := db.Query("SELECT id, username, title FROM chats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		chat := Chat{}
		err := rows.Scan(&(chat.ID), &(chat.Username), &(chat.Title))
		if err != nil {
			log.Fatal(err)
		}
		chats[chat.ID] = chat
	}
}
